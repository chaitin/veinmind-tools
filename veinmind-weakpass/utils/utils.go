package utils

import (
	"bufio"
	"errors"
	"fmt"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/dict"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/model"
	"os"
	"path/filepath"
	"strings"
)

func GetImageName(image api.Image) (imageName string, err error) {
	repoRefs, err := image.RepoRefs()
	if err != nil {
		return "unknow", err
	}
	if len(repoRefs) >= 1 {
		repoRefSplit := strings.Split(repoRefs[0], "/")
		imageName = repoRefSplit[len(repoRefSplit)-1]
		imageNameSplit := strings.Split(imageName, "@")
		imageName = imageNameSplit[0]
		imageNameSplit = strings.Split(imageName, ":")
		imageName = imageNameSplit[0]
	} else {
		imageName = image.ID()
	}
	return imageName, nil
}

func GetConfig(c *cmd.Command) model.Config {
	opt := model.Config{
		Thread: func() int {
			threads, err := c.Flags().GetInt("threads")
			if err != nil {
				return 10
			} else {
				return threads
			}
		}(),
		Username: func() string {
			username, err := c.Flags().GetString("username")
			if err != nil {
				return ""
			} else {
				return username
			}
		}(),
		Dictpath: func() string {
			dictpath, err := c.Flags().GetString("dictpath")
			if err != nil {
				return ""
			} else {
				return dictpath
			}
		}(),
	}
	return opt
}

func StartModule(conf model.Config, image api.Image, modulename string) (results model.ScanImageResult, err error) {

	// 初始化一个镜像扫描结果
	result := model.ScanImageResult{}
	result.ImageName, err = GetImageName(image)
	result.ImageID = image.ID()

	// 创建一个扫描模块
	mod, err := model.GetModuleByName(modulename)
	if err != nil {
		return model.ScanImageResult{}, err
	}

	log.Info(fmt.Sprintf("scan %s weakpass in %s", mod.Name(), result.ImageName))

	// 从配置文件中扩展字典
	var finalDict = []string{}
	baseDict := dict.Passdict
	if conf.Dictpath != "" {
		f, err := os.Open(conf.Dictpath)
		if err != nil {
			return model.ScanImageResult{}, errors.New("Dictpath not found")
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			baseDict = append(baseDict, scanner.Text())
		}
	}
	// 处理密码字典中的宏
	imageName, err := GetImageName(image)
	if err != nil {
		return model.ScanImageResult{}, err
	}
	for _, guess := range baseDict {
		// 动态替换弱口令字典中的宏
		if imageName != "" {
			guess = strings.Replace(guess, "${image_name}", imageName, -1)
		}
		finalDict = append(finalDict, guess)
	}

	// 从提供的默认路径中爆破
	pathes := mod.GetFilePath()
	for _, path := range pathes {
		var tmp = model.PasswdInfo{}
		file, err := image.Open(path)
		if err != nil {
			log.Warn(fmt.Sprintf("%s config file doesn't exist in image: %s!", path, result.ImageName))
			continue
		}
		err = mod.Init(conf)
		if err != nil {
			log.Error(err)
			continue
		}
		var Passwdinfos []model.PasswdInfo = []model.PasswdInfo{}
		PasswdinfosTmp, err := mod.ParsePasswdInfo(file)
		if err != nil {
			log.Warn(fmt.Sprintf("%s format error!", filepath.Base(path)))
			continue
		}
		for _, i := range PasswdinfosTmp {
			tmp.Username = i.Username
			tmp.Password = i.Password
			tmp.Filepath = path
			Passwdinfos = append(Passwdinfos, tmp)
		}
		// 开始密码爆破
		weakpass, err := mod.BrutePasswd(conf, Passwdinfos, finalDict, mod.MatchPasswd)

		if err != nil {
			log.Warn(err)
			continue
		}
		// 将扫描的所有弱密码,放到镜像扫描结果中
		result.WeakpassResults = append(result.WeakpassResults, weakpass...)
	}
	if len(result.WeakpassResults) > 0 {
		return result, nil
	} else {
		//镜像中没有该项服务
		return model.ScanImageResult{}, errors.New(fmt.Sprintf("%s doesn't exists in %s", mod.Name(), result.ImageName))
	}

}

func Findpathes(image api.Image, path string) (pathes []string) {
	filename := filepath.Base(path)
	image.Walk("/", func(path string, info os.FileInfo, err error) error {
		if info.Name() == filename {
			pathes = append(pathes, path)
		}
		return nil
	})
	return pathes
}
