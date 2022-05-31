package module

import (
	"errors"
	"github.com/Jeffail/tunny"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/veinmind-common/go/service/report"
	"io"
	"sync"
	"strings"
	"time"
)

type IModule interface {
	// 模块初始化
	Init(conf Config) error

	// 获取模块名称
	Name() string

	// 从文件中解析密码相关信息
	ParsePasswdInfo(file io.Reader) ([]PasswdInfo, error)

	// 密码匹配模式
	MatchPasswd(passwd string, guess string) bool

	// 爆破密码
	BrutePasswd(PasswdInfos []PasswdInfo, dicts []string, fn func(password, guess string) bool) ([]WeakpassResult, error)

	// 获取模块默认路径
	GetFilePath() []string

	// 获取特定的字典
	GetSpecialDict() []string

	// 根据模块名称处理字典
	ProcessDict(dict []string, moduleName string)

	// 生成报告
	GenerateReport(weakpassResults []WeakpassResult) (report.ReportEvent, error)
}

type Module struct {
	conf       Config
	name       string
	filePath   []string
	specialDict []string
	passwdType PasswordType
}

func (this *Module) Name() string {
	return this.name
}

func (this *Module) GetSpecialDict() []string {
	return this.specialDict
}

func (this *Module) GetFilePath() []string {
	return this.filePath
}

func (this *Module) ProcessDict(dict []string, moduleName string){
	for i, guess := range dict {
		// 根据模块名动态替换弱口令字典中的宏
		if moduleName != "" {
			dict[i] = strings.Replace(guess, "${module_name}", moduleName, -1)
		}
	}
}

func (this *Module) Init(config Config) (err error) {
	this.conf = config
	return nil
}

func (this *Module) ParsePasswdInfo(file io.Reader) ([]PasswdInfo, error) {
	panic("Module.ParsePasswdInfo() not implemented yet")
}

func (this *Module) BrutePasswd(PasswdInfos []PasswdInfo, dicts []string, fn func(password, guess string) bool) (weakpassResults []WeakpassResult, err error) {
	config := this.conf
	var weakpassResultsLock sync.Mutex
	// initial the concurrency pool
	pool := tunny.NewFunc(config.Thread, func(opt interface{}) interface{} {
		bruteOpt, ok := opt.(BruteOption)
		if !ok {
			return errors.New("please use BruteOption")
		}
		match := fn(bruteOpt.Passwdinfo.Password, bruteOpt.Guess)
		if match {
			w := WeakpassResult{
				Username: bruteOpt.Passwdinfo.Username,
				Password: bruteOpt.Guess,
				Filepath: bruteOpt.Passwdinfo.Filepath,
			}
			weakpassResultsLock.Lock()
			weakpassResults = append(weakpassResults, w)
			weakpassResultsLock.Unlock()
			return true
		}
		return false
	})
	defer pool.Close()

	for _, s := range PasswdInfos {
		// 判断是否为指定用户名
		if config.Username != "" {
			if s.Username != config.Username {
				continue
			}
		}
		for _, guess := range dicts {
			match, err := pool.ProcessTimed(BruteOption{
				Passwdinfo: s,
				Guess:      guess,
			}, 5*time.Second)

			if err != nil {
				return []WeakpassResult{}, err
			}

			if v, ok := match.(bool); ok {
				if v {
					break
				}
			}
		}
	}
	if len(weakpassResults) > 0{
		_, err = this.GenerateReport(weakpassResults)
		if err != nil {
			log.Warn("Report failed! cause ",err)
		}
	}
	return weakpassResults, nil

}

func (this *Module) MatchPasswd(passwd string, guess string) bool {
	return passwd == guess
}

func (this *Module) GenerateReport(weakpassResults []WeakpassResult) (Reportevent report.ReportEvent, err error) {
	details := []report.AlertDetail{}
	for _, wr := range weakpassResults {
		details = append(details, report.AlertDetail{
			WeakpassDetail: &report.WeakpassDetail{
				Username: wr.Username,
				Password: wr.Password,
				Service:  report.WeakpassService(this.passwdType),
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
			return report.ReportEvent{}, err
		}
		return Reportevent, nil
	}
	return report.ReportEvent{}, errors.New("Report is empty!")
}
