package malicious

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"debug/elf"
	"encoding/hex"
	"io"
	"io/fs"
	"net"
	"os"
	"strings"
	"syscall"
	"time"

	"code.cloudfoundry.org/bytefmt"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/libveinmind/go/plugin/log"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/database"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/database/model"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/sdk/av"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/sdk/av/clamav"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/sdk/av/virustotal"
)

func Scan(image api.Image) (scanReport model.ReportImage, err error) {
	// 判断是否已经扫描过
	database.GetDbInstance().Preload("Layers").Preload("Layers.MaliciousFileInfos").Where("image_id = ?", image.ID()).Find(&scanReport)
	if scanReport.ImageID != "" {
		log.Info(image.ID(), " Has been detected")
		return scanReport, nil
	}

	refs, err := image.RepoRefs()
	var imageRef string
	if err == nil && len(refs) > 0 {
		imageRef = refs[0]
	} else {
		imageRef = image.ID()
	}
	log.Info("Scan Image: ", imageRef)

	// 判断是否可以获取 Layer
	switch v := image.(type) {
	case *docker.Image:
		dockerImage := v
		for i := 0; i < dockerImage.NumLayers(); i++ {
			// 获取 Layer ID
			layerID, err := dockerImage.GetLayerDiffID(i)
			if err != nil {
				log.Error("Get LayerID Error: ", err)
				continue
			}

			// 判断 Layer 是否已经扫描
			reportLayer := model.ReportLayer{}
			database.GetDbInstance().Preload("MaliciousFileInfos").Where("layer_id", layerID).Find(&reportLayer)

			if reportLayer.LayerID != "" {
				reportLayerCopy := model.ReportLayer{
					ImageID:            image.ID(),
					LayerID:            reportLayer.LayerID,
					MaliciousFileInfos: reportLayer.MaliciousFileInfos,
				}
				scanReport.Layers = append(scanReport.Layers, reportLayerCopy)
				log.Info("Skip Scan Layer: ", layerID)
				continue
			} else {
				l, err := dockerImage.OpenLayer(i)
				if err != nil {
					log.Error(err)
				}

				log.Info("Start Scan Layer: ", l.ID())
				l.Walk("/", func(path string, info fs.FileInfo, err error) error {
					// 部分情况下ELF解析会产生panic
					defer func() {
						if err := recover(); err != nil {
							log.Error(err)
						}
					}()

					// 处理错误
					if err != nil {
						log.Debug(err)
						return nil
					}

					// 判断文件类型，跳过特定类型文件
					if (info.Mode() & (os.ModeDevice | os.ModeNamedPipe | os.ModeSocket | os.ModeCharDevice | os.ModeDir)) != 0 {
						log.Debug("Skip: ", path)
						return nil
					}

					// 忽略软链接, PS: 全局扫描终究会扫到实际的文件
					if (info.Mode() & os.ModeSymlink) != 0 {
						log.Debug("Skip: ", path)
						return nil
					}

					scanReport.ScanFileCount++

					f, err := l.Open(path)
					if err != nil {
						log.Debug(err)
						return nil
					}

					defer func() {
						f.Close()
					}()

					// 判断是否是ELF文件，如果不是则跳过
					_, err = elf.NewFile(f)
					if _, ok := err.(*elf.FormatError); ok {
						log.Debug("Skip File: ", path)
						return nil
					} else if err != nil {
						return nil
					}

					var results []av.ScanResult

					// 使用 ClamAV 进行扫描
					if clamav.Active() {
						results, err = clamav.ScanStream(f)
						if err != nil {
							if _, ok := err.(*net.OpError); ok {
								log.Error(err)
							} else {
								//TODO: 告知使用者其他Err
								log.Debug(err)
							}
						}
					}

					// 使用 Virustotal 进行扫描
					fileByte, err := io.ReadAll(f)
					hash := sha256.New()
					fileSha256 := hex.EncodeToString(hash.Sum(fileByte))

					virustotalContext, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
					if virustotal.Active() {
						vtResults, err := virustotal.ScanSHA256(virustotalContext, fileSha256)
						if err == nil && vtResults != nil && len(vtResults) > 0 {
							results = append(results, vtResults...)
						}
					}

					if len(results) > 0 {
						log.Warn("Find malicious file: ", path)

						// 假设有多个结果，直接拼接 description
						description := ""
						engine := map[string]bool{}
						engineName := ""
						for _, r := range results {
							description = description + r.Description + ","
							engine[r.EngineName] = true
						}
						for e := range engine {
							engineName = e + ","
						}
						engineName = strings.TrimRight(engineName, ",")
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
							Engine:       engineName,
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
						log.Warn(r)
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
		log.Error(err)
	}

	repoRefs, err := image.RepoRefs()
	if err == nil && len(repoRefs) >= 1 {
		scanReport.ImageName = repoRefs[0]
	} else {
		scanReport.ImageName = image.ID()
	}

	// 存储结果
	log.Info("Store Scan Report: ", image.ID())
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
