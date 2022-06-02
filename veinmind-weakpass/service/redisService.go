package service

import (
	"io"

	"github.com/chaitin/veinmind-tools/veinmind-weakpass/extractor"
)

type redisService struct {
	name     string
	filepath []string
	extractor.Extractor
}

func (i *redisService) Name() string {
	return "redis"
}

func (i *redisService) FilePath() (paths []string) {
	return i.filepath
}
func (i *redisService) GetRecords(file io.Reader) (records []extractor.Record, err error) {
	Extractor, err := GetExtractor("redis")
	if err != nil {
		return records, err
	}
	// 从文件中获取密码相关的记录
	records, err = Extractor.Extract(file)
	if err != nil {
		return records, err
	}
	return records, nil
}

func init() {
	mod := &redisService{}
	mod.name = "redis"
	mod.filepath = []string{"/etc/redis/redis.conf", "/etc/redisc.conf"}
	Register(mod)
}
