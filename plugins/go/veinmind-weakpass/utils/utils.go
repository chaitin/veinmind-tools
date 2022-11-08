package utils

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Jeffail/tunny"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/model"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/service"
)

func GetImageName(image api.Image) (imageName string, err error) {
	repoRefs, err := image.RepoRefs()
	if err != nil {
		return "unknown", err
	}
	if len(repoRefs) >= 1 {
		imageName = repoRefs[0]
	} else {
		imageName = image.ID()
	}
	return imageName, nil
}

func StartModule(config model.Config, fs api.FileSystem, modname string, marco map[string]string) (results []model.WeakpassResult, err error) {
	// 获取对应的服务模块
	mod, err := service.GetModuleByName(modname)
	if err != nil {
		return results, err
	}
	// 获取对应模块的加密算法
	hash, err := service.GetHash(modname)
	if err != nil {
		return results, err
	}

	// 最终需要爆破的字典
	finalDict := []string{}
	// 加载模块对应的字典
	finalDict = append(finalDict, service.GetDict(modname)...)
	// 读取用户提供的字典
	if config.Dictpath != "" {
		f, err := os.Open(config.Dictpath)
		if err != nil {
			// 如果用户配置的字典打开失败,日志告知即可,默认字典依然可以使用
			log.Warn("use default dict cause", err)
		}
		if err == nil {
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				finalDict = append(finalDict, scanner.Text())
			}
			f.Close()
		}
	}

	// 获取对应模块的filepaths
	filepaths := mod.FilePath()
	// 开始从路径中爆破密码
	WeakpassResults := []model.WeakpassResult{}
	for _, path := range filepaths {
		_, err := fs.Stat(path)
		if err != nil {
			// 如果镜像中不存在配置文件,直接跳过
			continue
		}
		file, err := fs.Open(path)
		if err != nil {
			log.Warn(err)
			continue
		}
		// 创建线程池,加快爆破速度
		var weakpassResultsLock sync.Mutex
		pool := tunny.NewFunc(config.Thread, func(opt interface{}) interface{} {
			bruteOpt, ok := opt.(model.BruteOption)
			if !ok {
				return errors.New("please use BruteOption")
			}
			match, err := hash.Match(bruteOpt.Records.Password, bruteOpt.Guess)
			if err != nil {
				return err
			}
			if match {
				w := model.WeakpassResult{
					Username:    bruteOpt.Records.Username,
					Password:    bruteOpt.Guess,
					Filepath:    path,
					ServiceType: service.GetType(mod),
				}
				weakpassResultsLock.Lock()
				WeakpassResults = append(WeakpassResults, w)
				weakpassResultsLock.Unlock()
				return true
			}
			return false
		})
		defer pool.Close()

		// 从文件中获取密码相关的记录
		records, err := mod.GetRecords(file)
		file.Close()
		for _, item := range records {
			// 判断是否为指定用户名
			if config.Username != "" {
				if item.Username != config.Username {
					continue
				}
			}
			for _, guess := range finalDict {
				// 替换镜像名相关的宏
				guess = strings.Replace(guess, "${image_name}", marco["image_name"], -1)
				// 替换服务名相关的宏
				guess = strings.Replace(guess, "${module_name}", marco["module_name"], -1)

				match, err := pool.ProcessTimed(model.BruteOption{
					Records: item,
					Guess:   guess,
				}, 5*time.Second)

				if err != nil {
					log.Error(err)
					continue
				}
				if v, ok := match.(bool); ok {
					if v {
						break
					}
				}
			}
		}

	}

	if len(WeakpassResults) > 0 {
		results = append(results, WeakpassResults...)
	}

	return results, nil

}

func GenerateImageReport(weakpassResults []model.WeakpassResult, image api.Image) (err error) {
	details := []report.AlertDetail{}
	for _, wr := range weakpassResults {
		details = append(details, report.AlertDetail{
			WeakpassDetail: &report.WeakpassDetail{
				Username: wr.Username,
				Password: wr.Password,
				Service:  wr.ServiceType,
				Path:     wr.Filepath,
			},
		})
	}
	if len(details) > 0 {
		Reportevent := report.ReportEvent{
			ID:           image.ID(),
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

func GenerateContainerReport(weakpassResults []model.WeakpassResult, container api.Container) (err error) {
	details := []report.AlertDetail{}
	for _, wr := range weakpassResults {
		details = append(details, report.AlertDetail{
			WeakpassDetail: &report.WeakpassDetail{
				Username: wr.Username,
				Password: wr.Password,
				Service:  wr.ServiceType,
				Path:     wr.Filepath,
			},
		})
	}
	if len(details) > 0 {
		Reportevent := report.ReportEvent{
			ID:           container.ID(),
			Time:         time.Now(),
			Level:        report.High,
			DetectType:   report.Container,
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
