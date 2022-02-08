package scan

import (
	"bufio"
	"errors"
	"github.com/Jeffail/tunny"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/containerd"
	"github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/brute"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/brute/ssh_passwd"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/embed"
	common "github.com/chaitin/veinmind-tools/veinmind-weakpass/log"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/model"
	"os"
	"strings"
	"sync"
	"time"
)

var passDict []string

func init() {
	// 初始化字典
	passDictFile, err := embed.EmbedFS.Open("pass.dict")
	if err != nil {
		common.Log.Fatal(err)
	}

	scanner := bufio.NewScanner(passDictFile)
	for scanner.Scan() {
		passDict = append(passDict, scanner.Text())
	}
}

type SSHScanPlugin struct {
}

type SSHBruteOpt struct {
	Shadow ssh_passwd.Shadow
	Guess  string
}

func (self *SSHScanPlugin) Scan(opt ScanOption) (results []model.ScanImageResult, err error) {
	// 初始化客户端
	var client api.Runtime

	switch opt.EngineType {
	case Dockerd:
		client, err = docker.New()
		if err != nil {
			return nil, err
		}

		defer func() {
			client.Close()
		}()
	case Containerd:
		client, err = containerd.New()

		if err != nil {
			return nil, err
		}

		defer func() {
			client.Close()
		}()
	default:
		return nil, errors.New("Engine type doesn't match")
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
		scanResult, err := self.ScanById(imageID, client, opt)
		if err != nil {
			common.Log.Error(err)
			continue
		}

		results = append(results, scanResult)
	}

	return results, nil
}

func (self *SSHScanPlugin) ScanById(id string, client api.Runtime, opt ScanOption) (model.ScanImageResult, error) {
	image, err := client.OpenImageByID(id)
	if err != nil {
		return model.ScanImageResult{}, err
	}

	f, err := image.Open("/etc/shadow")
	if err != nil {
		return model.ScanImageResult{}, err
	}

	// 设置镜像报告信息
	imageResult := model.ScanImageResult{}

	// 设置镜像名称
	repoRefs, err := image.RepoRefs()
	if err == nil && len(repoRefs) >= 1 {
		imageResult.ImageName = repoRefs[0]
	} else {
		imageResult.ImageName = image.ID()
	}
	imageResult.ImageID = image.ID()

	// 获取镜像名称(排除仓库地址和namespace)
	var imageName string
	if len(repoRefs) >= 1 {
		repoRefSplit := strings.Split(repoRefs[0], "/")
		imageName = repoRefSplit[len(repoRefSplit)-1]
		imageNameSplit := strings.Split(imageName, "@")
		imageName = imageNameSplit[0]
		imageNameSplit = strings.Split(imageName, ":")
		imageName = imageNameSplit[0]
	}

	// 解析shadow配置文件
	shadows, err := ssh_passwd.ParseShadowFile(f)

	// 检测结果
	var weakpassResultsLock sync.Mutex
	var weakpassResults []model.WeakpassResult

	// 初始化并发池
	pool := tunny.NewFunc(opt.ScanThreads, func(opt interface{}) interface{} {
		bruteOpt, ok := opt.(SSHBruteOpt)
		if !ok {
			return errors.New("Please use sshbruteopt")
		}

		_, match := brute.SSHMatchPassword(bruteOpt.Shadow.EncryptedPassword, bruteOpt.Guess)
		if match {
			w := model.WeakpassResult{
				PassType: model.SSH,
				Username: bruteOpt.Shadow.LoginName,
				Password: bruteOpt.Guess,
			}

			weakpassResultsLock.Lock()
			weakpassResults = append(weakpassResults, w)
			weakpassResultsLock.Unlock()

			return true
		}

		return false
	})
	defer pool.Close()

	// 初始化字典
	if opt.Dictpath != "" {
		f, err := os.Open(opt.Dictpath)
		if err != nil {
			common.Log.Fatal(err)
		}
		scanner := bufio.NewScanner(f)
		passDict = []string{}
		for scanner.Scan() {
			passDict = append(passDict, scanner.Text())
		}
	}

	common.Log.Info("Start Scan Image SSH Weakpass: ", imageResult.ImageName)

	for _, s := range shadows {
		// 判断是否为指定用户名
		if opt.Username != "" {
			if s.LoginName != opt.Username {
				continue
			}
		}

		for _, guess := range passDict {
			// 动态替换弱口令字典中的宏
			if imageName != "" {
				guess = strings.Replace(guess, "${image_name}", imageName, -1)
			}

			match, err := pool.ProcessTimed(SSHBruteOpt{
				Shadow: s,
				Guess:  guess,
			}, 5*time.Second)

			if err != nil {
				common.Log.Error(err)
			}

			if v, ok := match.(bool); ok {
				if v {
					break
				}
			}
		}
	}

	imageResult.WeakpassResults = weakpassResults

	return imageResult, nil
}
