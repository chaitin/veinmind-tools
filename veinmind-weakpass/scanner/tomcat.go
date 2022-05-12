package scanner

import (
	"bufio"
	"errors"
	"github.com/Jeffail/tunny"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/embed"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/brute"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/model"
	"os"
	"strings"
	"sync"
	"time"
	"fmt"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/brute/tomcat_passwd"
)


func init() {
	// 初始化字典
	passDictFile, err := embed.EmbedFS.Open("pass.dict")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(passDictFile)
	for scanner.Scan() {
		passDict = append(passDict, scanner.Text())
	}
}


type TomcatBruteOpt struct {
	Tomcat tomcat_passwd.Tomcat
	Guess string
}

func mergeArr(a, b []tomcat_passwd.Tomcat,c string) []tomcat_passwd.Tomcat {
	var arr []tomcat_passwd.Tomcat
	for _, i := range a {
	   arr = append(arr, i)
	}
	for _, j := range b {
		j.Filepath = c
	   arr = append(arr, j)
	}
	return arr
 }
 
func ScanTomcat(image api.Image, opt ScanOption) (model.ScanImageResult, error) {
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
	// 寻找config_user.xml,会不会有两个这样的文件,如果有是否都需要遍历
	var tomcat_passwd_files []string 
	image.Walk("/", func (path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("ERROR: %v", err)
			return err
		}
		if info.Name() == "tomcat-users.xml" {
			parent_len := len(path)-16-5
			//fmt.Printf(path[:parent_len] + "bin/catalina-tasks.xml")
			info,err = image.Stat(path[:parent_len]+"bin/catalina-tasks.xml")
			if err == nil{
				tomcat_passwd_files = append(tomcat_passwd_files,path)
			}
		}
		return nil
	})
	var tomcats []tomcat_passwd.Tomcat
	for _,tomcat_passwd_file := range tomcat_passwd_files{
		file, err := image.Open(tomcat_passwd_file)
		if err != nil {
			return model.ScanImageResult{}, err
		}
		tomcat,err := tomcat_passwd.ParseTomcatFile(file)
		tomcats = mergeArr(tomcats,tomcat,tomcat_passwd_file)
	}	
	
	// 检测结果
	var weakpassResultsLock sync.Mutex
	var weakpassResults []model.WeakpassResult

	// 初始化并发池
	pool := tunny.NewFunc(opt.ScanThreads, func(opt interface{}) interface{} {
		bruteOpt, ok := opt.(TomcatBruteOpt)
		if !ok {
			return errors.New("please use tomcatbruteopt")
		}

		_, match := brute.TomcatMatchPassword(bruteOpt.Tomcat.Username, bruteOpt.Guess)
		if match {
			w := model.WeakpassResult{
				PassType: model.TOMCAT,
				Username: bruteOpt.Tomcat.Username,
				Password: bruteOpt.Guess,
				Filepath: bruteOpt.Tomcat.Filepath,
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
			log.Fatal(err)
		}
		scanner := bufio.NewScanner(f)
		passDict = []string{}
		for scanner.Scan() {
			passDict = append(passDict, scanner.Text())
		}
	}

	log.Info("start scan image tomcat weakpass: ", imageResult.ImageName)

	for _, s := range tomcats {
		// 判断是否为指定用户名
		if opt.Username != "" {
			if s.Username != opt.Username {
				continue
			}
		}

		for _, guess := range passDict {
			// 动态替换弱口令字典中的宏
			if imageName != "" {
				guess = strings.Replace(guess, "${image_name}", imageName, -1)
			}

			match, err := pool.ProcessTimed(TomcatBruteOpt{
				Tomcat: s,
				Guess:  guess,
			}, 5*time.Second)

			if err != nil {
				log.Error(err)
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