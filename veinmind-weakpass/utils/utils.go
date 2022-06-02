package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/veinmind-common/go/service/report"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/model"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/service"
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

func StartModule(c *cmd.Command, image api.Image, modname string) (result model.ScanImageResult, err error) {
	// 从cmd中获取相关的配置信息
	config := GetConfig(c)
	// 初始化一个镜像扫描结果
	result = model.ScanImageResult{}
	result.ServiceName = modname
	imagename, err := GetImageName(image)
	if err != nil {
		// 告知镜像名称有问题,使用unkonow
		log.Warn(err)
	}
	result.ImageName = imagename
	result.ImageID = image.ID()

	// 获取对应的服务模块
	mod, err := service.GetModuleByName(modname)
	if err != nil {
		log.Warn(err)
		return result, err
	}

	// 最终需要爆破的字典
	finalDict := []string{}
	// 加载模块对应的字典
	finalDict = append(finalDict, service.GetDict(modname)...)
	// 读取本地的字典
	if config.Dictpath != "" {
		f, err := os.Open(config.Dictpath)
		if err != nil {
			// 如果用户配置的字典打开失败,日志告知即可,默认字典依然可以使用
			log.Warn("use default dict cause", err)
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			finalDict = append(finalDict, scanner.Text())
		}
	}
	// 替换字典中的宏
	for i, guess := range finalDict {
		// 替换镜像名相关的宏
		if imagename != "" {
			finalDict[i] = strings.Replace(guess, "${image_name}", imagename, -1)
		}
		// 替换服务名相关的宏
		finalDict[i] = strings.Replace(guess, "${module_name}", modname, -1)
	}

	// 获取对应模块的filepaths
	filepaths := mod.FilePath()
	log.Info(fmt.Sprintf("start to scan %s in image %s", modname, imagename))
	// 开始从路径中爆破密码
	WeakpassResults := []model.WeakpassResult{}
	for _, path := range filepaths {
		_, err := image.Stat(path)
		if err != nil {
			// 如果镜像中不存在配置文件,直接跳过
			continue
		}
		file, err := image.Open(path)
		if err != nil {
			log.Warn(err)
			continue
		}
		// 从文件中获取密码相关的记录
		records, err := mod.GetRecords(file)
		for _, item := range records {
			plain, err := item.Password.Match(finalDict)
			if err != nil {
				// 密码不匹配
				continue
			}
			var tmp = model.WeakpassResult{Username: item.Username, Password: plain, Filepath: path}
			WeakpassResults = append(WeakpassResults, tmp)
		}

	}

	// 进行Report
	if len(WeakpassResults) > 0 {
		err = GenerateReport(WeakpassResults)
		if err != nil {
			log.Warn(err)
		}
		result.WeakpassResults = WeakpassResults
	}

	return result, nil

}

func GenerateReport(weakpassResults []model.WeakpassResult) (err error) {
	details := []report.AlertDetail{}
	for _, wr := range weakpassResults {
		details = append(details, report.AlertDetail{
			WeakpassDetail: &report.WeakpassDetail{
				Username: wr.Username,
				Password: wr.Password,
			},
		})
	}
	if len(details) > 0 {
		Reportevent := report.ReportEvent{
			Time:         time.Now(),
			Level:        report.High,
			DetectType:   report.Image,
			EventType:    report.Risk,
			AlertType:    report.Weakpass,
			AlertDetails: details,
		}
		err = report.DefaultReportClient().Report(Reportevent)
		if err != nil {
			return err
		}
	}
	return nil
}
