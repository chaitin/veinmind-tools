package scanner

import (
	"bufio"
	"errors"
	"github.com/Jeffail/tunny"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/brute"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/brute/redis"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/embed"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/model"
	"os"
	"strings"
	"sync"
	"time"
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

type RedisBruteOpt struct {
	Redis redis.Redis
	Guess string
}

func ScanRedis(image api.Image, opt ScanOption) (model.ScanImageResult, error) {
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
	var RedisPasswdFiles = []string{"/etc/redis.conf", "/etc/redis/redis.conf"}
	// image.Walk("/", func(path string, info os.FileInfo, err error) error {
	// 	if info.Name() == "redis.conf" {
	// 		RedisPasswdFiles = append(RedisPasswdFiles, path)
	// 	}
	// 	return nil
	// })
	var rediss []redis.Redis
	for _, RedisPasswdFile := range RedisPasswdFiles {
		_, err := image.Stat(RedisPasswdFile)
		if err != nil {
			continue
		}
		log.Error(RedisPasswdFile)
		file, err := image.Open(RedisPasswdFile)
		if err != nil {
			return model.ScanImageResult{}, err
		}
		redis, err := redis.ParseRedisFile(file)
		if err != nil {
			return model.ScanImageResult{}, err
		}
		for _, i := range redis {
			i.Filepath = RedisPasswdFile
		}
		rediss = append(rediss, redis...)
	}

	// 检测结果
	var weakpassResultsLock sync.Mutex
	var weakpassResults []model.WeakpassResult

	// 初始化并发池
	pool := tunny.NewFunc(opt.ScanThreads, func(opt interface{}) interface{} {
		bruteOpt, ok := opt.(TomcatBruteOpt)
		if !ok {
			return errors.New("please use redisbruteopt")
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

	log.Info("start scan image redis weakpass: ", imageResult.ImageName)

	for _, s := range rediss {
		for _, guess := range passDict {
			// 动态替换弱口令字典中的宏
			if imageName != "" {
				guess = strings.Replace(guess, "${image_name}", imageName, -1)
			}

			match, err := pool.ProcessTimed(RedisBruteOpt{
				Redis: s,
				Guess: guess,
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
