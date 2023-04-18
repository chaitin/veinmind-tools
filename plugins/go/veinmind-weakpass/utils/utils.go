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
	"github.com/chaitin/veinmind-common-go/service/report/event"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/model"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/service"
)

func StartModule(config model.Config, fs api.FileSystem, modname string, marco map[string]string) (results []model.WeakpassResult, err error) {
	// 获取对应的服务模块
	mods, err := service.GetModuleByName(modname)
	if err != nil {
		return results, err
	}
	// 遍历对应模块
	for _, mod := range mods {

		// 获取对应模块的加密算法
		hash, err := service.GetHash(mod.Name())
		if err != nil {
			return results, err
		}

		// 最终需要爆破的字典
		var finalDict []string
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
				err := f.Close()
				if err != nil {
					return nil, err
				}
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
			err = file.Close()
			if err != nil {
				return nil, err
			}
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
	}

	return results, nil

}

func GenerateImageReport(weakpassResults []model.WeakpassResult, image api.Image, reportService *report.Service) (err error) {
	var details []event.AlertDetail
	for _, wr := range weakpassResults {
		details = append(details, &event.WeakpassDetail{
			Username: wr.Username,
			Password: wr.Password,
			Service:  wr.ServiceType,
			Path:     wr.Filepath,
		})
	}
	for _, d := range details {
		event := &event.Event{
			BasicInfo: &event.BasicInfo{
				ID:         image.ID(),
				Object:     event.NewObject(image),
				Source:     "veinmind-weakpass",
				Time:       time.Now(),
				Level:      event.High,
				DetectType: event.Image,
				EventType:  event.Risk,
				AlertType:  event.Weakpass,
			},
			DetailInfo: &event.DetailInfo{
				AlertDetail: d,
			},
		}
		err = reportService.Client.Report(event)
	}
	if err != nil {
		return err
	}
	return nil
}

func GenerateContainerReport(weakpassResults []model.WeakpassResult, container api.Container, reportService *report.Service) (err error) {
	var details []event.AlertDetail
	for _, wr := range weakpassResults {
		details = append(details, &event.WeakpassDetail{
			Username: wr.Username,
			Password: wr.Password,
			Service:  wr.ServiceType,
			Path:     wr.Filepath,
		})
	}
	for _, d := range details {
		event := &event.Event{
			BasicInfo: &event.BasicInfo{
				ID:         container.ID(),
				Object:     event.NewObject(container),
				Source:     "veinmind-weakpass",
				Time:       time.Now(),
				Level:      event.High,
				DetectType: event.Container,
				EventType:  event.Risk,
				AlertType:  event.Weakpass,
			},
			DetailInfo: &event.DetailInfo{
				AlertDetail: d,
			},
		}
		err = reportService.Client.Report(event)
	}
	if err != nil {
		return err
	}
	return nil
}
