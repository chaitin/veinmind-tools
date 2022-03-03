package malicious

import (
	"code.cloudfoundry.org/bytefmt"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"debug/elf"
	"encoding/hex"
	"errors"
	veinmindcommon "github.com/chaitin/libveinmind/go"
	containerd "github.com/chaitin/libveinmind/go/containerd"
	docker "github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/database"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/database/model"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/scanner/scanner_common"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/sdk/av"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/sdk/av/clamav"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/sdk/av/virustotal"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/sdk/common"
	fs "io/fs"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"syscall"
	"time"
)

type MaliciousPlugin struct {
}

func (self *MaliciousPlugin) Scan(opt scanner_common.ScanOption) (scanReportAll model.ReportData, err error) {
	// 扫描之前需要至少拥有一个存活的 AV 引擎， 否则抛出异常
	if !clamav.Active() && !virustotal.Active() {
		log.Fatal("Please active anti virus engine at least one")
	}

	// 判断引擎类型
	var client veinmindcommon.Runtime

	switch opt.EngineType {
	case scanner_common.Dockerd:
		dockerClient, err := docker.New()
		if err != nil {
			return model.ReportData{}, err
		}

		client = dockerClient

		defer func() {
			client.Close()
		}()
	case scanner_common.Containerd:
		containerClient, err := containerd.New()
		if err != nil {
			return model.ReportData{}, err
		}

		client = containerClient

		defer func() {
			client.Close()
		}()
	default:
		return model.ReportData{}, errors.New("Engine Type Not Match")
	}

	var imageIds []string
	if opt.ImageName != "" {
		imageIds, err = client.FindImageIDs(opt.ImageName)
		if err != nil {
			return
		}
	} else {
		imageIds, err = client.ListImageIDs()
		if err != nil {
			return
		}
	}

	for _, imageID := range imageIds {
		scanResult, err := self.ScanById(imageID, client)
		if err != nil {
			common.Log.Error(err)
			continue
		}

		scanReportAll.ScanImageResult = append(scanReportAll.ScanImageResult, scanResult)
	}

	return scanReportAll, nil
}

func (self *MaliciousPlugin) ScanById(id string, client veinmindcommon.Runtime) (scanReport model.ReportImage, err error) {
	// 判断是否已经扫描过
	database.GetDbInstance().Where("image_id = ?", id).Find(&scanReport)
	if scanReport.ImageID != "" {
		common.Log.Info(id, " Has been detected")
		return scanReport, nil
	}

	image, err := client.OpenImageByID(id)
	if err != nil {
		common.Log.Error(err)
		return
	}
	defer func() {
		image.Close()
	}()

	refs , err := image.RepoRefs()
	var imageRef string
	if err == nil && len(refs) > 0 {
		imageRef = refs[0]
	}else{
		imageRef = image.ID()
	}
	common.Log.Info("Scan Image: ", imageRef)

	// 判断是否可以获取 Layer
	switch v := image.(type) {
	case *docker.Image:
		dockerImage := v
		for i := 0; i < dockerImage.NumLayers(); i++ {
			// 获取 Layer ID
			layerID, err := dockerImage.GetLayerDiffID(i)
			if err != nil {
				common.Log.Error("Get LayerID Error: ", err)
				continue
			}

			// 判断 Layer 是否已经扫描
			reportLayer := model.ReportLayer{}
			database.GetDbInstance().Where("layer_id", layerID).Find(&reportLayer)

			if reportLayer.LayerID != "" {
				reportLayerCopy := model.ReportLayer{
					ImageID:            image.ID(),
					LayerID:            reportLayer.LayerID,
					MaliciousFileInfos: reportLayer.MaliciousFileInfos,
				}
				scanReport.Layers = append(scanReport.Layers, reportLayerCopy)
				common.Log.Info("Skip Scan Layer: ", layerID)
				continue
			} else {
				l, err := dockerImage.OpenLayer(i)
				if err != nil {
					common.Log.Error(err)
				}

				common.Log.Info("Start Scan Layer: ", l.ID())
				l.Walk("/", func(path string, info fs.FileInfo, err error) error {
					// 部分情况下ELF解析会产生panic
					defer func() {
						if err := recover(); err != nil {
							common.Log.Error(err)
						}
					}()

					// 处理错误
					if err != nil {
						common.Log.Debug(err)
						return nil
					}

					// 判断文件类型，跳过特定类型文件
					if (info.Mode() & (os.ModeDevice | os.ModeNamedPipe | os.ModeSocket | os.ModeCharDevice | os.ModeDir)) != 0 {
						common.Log.Debug("Skip: ", path)
						return nil
					}

					// 忽略软链接, PS: 全局扫描终究会扫到实际的文件
					if (info.Mode() & os.ModeSymlink) != 0 {
						common.Log.Debug("Skip: ", path)
						return nil
					}

					scanReport.ScanFileCount++

					f, err := l.Open(path)
					if err != nil {
						common.Log.Debug(err)
						return nil
					}

					defer func() {
						f.Close()
					}()

					// 判断是否是ELF文件，如果不是则跳过
					_, err = elf.NewFile(f)
					if _, ok := err.(*elf.FormatError); ok {
						common.Log.Debug("Skip File: ", path)
						return nil
					} else if err != nil {
						return nil
					}

					results := []av.ScanResult{}

					// 使用 ClamAV 进行扫描
					if clamav.Active() {
						results, err = clamav.ScanStream(f)
						if err != nil {
							if _, ok := err.(*net.OpError); ok {
								common.Log.Error(err)
							} else {
								//TODO: 告知使用者其他Err
								common.Log.Debug(err)
							}
						}
					}

					// 使用 Virustotal 进行扫描
					fileByte, err := ioutil.ReadAll(f)
					hash := sha256.New()
					fileSha256 := hex.EncodeToString(hash.Sum(fileByte))

					virustotalContext, _ := context.WithTimeout(context.Background(), 10 * time.Millisecond)
					if virustotal.Active() {
						vtResults, err := virustotal.ScanSHA256(virustotalContext, fileSha256)
						if err == nil && vtResults != nil && len(vtResults) > 0 {
							results = append(results, vtResults...)
						}
					}

					if len(results) > 0 {
						common.Log.Warn("Find malicious file: ", path)

						// 假设有多个结果，直接拼接 description
						description := ""
						for _, r := range results {
							description = description + r.Description + ","
						}
						description = strings.TrimRight(description, ",")

						scanReport.MaliciousFileCount++

						// 计算文件MD5
						hash := md5.New()
						var fileMd5 string
						if err == nil {
							fileMd5 = hex.EncodeToString(hash.Sum(fileByte)[:16])
						}

						// 获取文件时间
						stat := info.Sys().(*syscall.Stat_t)

						result := model.MaliciousFileInfo{
							RelativePath: path,
							FileName:     info.Name(),
							FileSize:     bytefmt.ByteSize(uint64(info.Size())),
							FileMd5:      fileMd5,
							FileSha256:   fileSha256,
							FileCreated:  time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec)).Format("2006-01-02 15:04:05"),
							Description:  description,
						}
						reportLayer.MaliciousFileInfos = append(reportLayer.MaliciousFileInfos, result)
					}
					for _, r := range results {
						common.Log.Warn(r)
					}
					return nil
				})

				reportLayer.LayerID = layerID
				reportLayer.ImageID = image.ID()
				scanReport.Layers = append(scanReport.Layers, reportLayer)
			}
		}
	}

	// 设置返回结果
	scanReport.ImageID = image.ID()
	oci, err := image.OCISpecV1()
	if err == nil && oci != nil {
		scanReport.ImageCreatedAt = oci.Created.Format("2006-01-02 15:04:05")
	} else {
		common.Log.Error(err)
	}

	repoRefs, err := image.RepoRefs()
	if err == nil && len(repoRefs) >= 1 {
		scanReport.ImageName = repoRefs[0]
	} else {
		scanReport.ImageName = image.ID()
	}

	// 存储结果
	common.Log.Info("Store Scan Report: ", image.ID())
	database.GetDbInstance().Create(&scanReport)
	for _, layerReport := range scanReport.Layers {
		l := layerReport
		database.GetDbInstance().Save(&l)
		for _, maliciousFile := range layerReport.MaliciousFileInfos {
			m := maliciousFile
			database.GetDbInstance().Save(&m)
		}
	}

	return scanReport, nil
}

func (self *MaliciousPlugin) PluginName() string {
	return "MaliciousPlugin"
}
