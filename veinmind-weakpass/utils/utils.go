package utils

import (
	"bufio"
	"errors"
	"fmt"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/dict"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/module"
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

func GetConfig(c *cmd.Command) module.Config {
	opt := module.Config{
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

func StartModule(conf module.Config, image api.Image, modulename string) (results module.ScanImageResult, err error) {

	// 初始化一个镜像扫描结果
	result := module.ScanImageResult{}
	result.ImageName, err = GetImageName(image)
	result.ImageID = image.ID()

	// 创建一个扫描模块
	mod, err := module.GetModuleByName(modulename)
	if err != nil {
		return module.ScanImageResult{}, err
	}

	log.Info(fmt.Sprintf("scan %s weakpass in %s", mod.Name(), result.ImageName))

	// 从配置文件中扩展字典
	var finalDict = []string{}
	// 加载模块对应的字典
	finalDict = append(finalDict,mod.GetSpecialDict()...)
	baseDict := dict.Passdict
	if conf.Dictpath != "" {
		f, err := os.Open(conf.Dictpath)
		if err != nil {
			return module.ScanImageResult{}, errors.New("Dictpath not found")
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			baseDict = append(baseDict, scanner.Text())
		}
	}
	finalDict = append(finalDict,baseDict...)

	// 根据模块名处理密码字典中的宏
	mod.ProcessDict(finalDict,modulename)

	// 根据镜像名称处理密码字典中的宏
	imageName, err := GetImageName(image)
	if err != nil {
		return module.ScanImageResult{}, err
	}
	for i, guess := range finalDict {
		// 动态替换弱口令字典中的宏
		if imageName != "" {
			finalDict[i] = strings.Replace(guess, "${image_name}", imageName, -1)
		}
	}

	// 从提供的默认路径中爆破
	pathes := mod.GetFilePath()
	for _, path := range pathes {
		var tmp = module.PasswdInfo{}
		_, err := os.Stat(path)
		if err != nil {
			// log.Warn(err)
			continue
		}
		file, err := image.Open(path)
		if err != nil {
			log.Warn(err)
			continue
		}
		err = mod.Init(conf)
		if err != nil {
			log.Error(err)
			continue
		}
		var Passwdinfos []module.PasswdInfo = []module.PasswdInfo{}
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
		weakpass, err := mod.BrutePasswd(Passwdinfos, finalDict, mod.MatchPasswd)

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
		return module.ScanImageResult{}, errors.New(fmt.Sprintf("%s doesn't exists in %s", mod.Name(), result.ImageName))
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
