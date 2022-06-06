package service

import (
	"io"
	"regexp"
	
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/model"
)

type redisService struct {
	name     string
	filepath []string
}

func (i *redisService) Name() string {
	return i.name
}

func (i *redisService) FilePath() (paths []string) {
	return i.filepath
}
func (i *redisService) GetRecords(file io.Reader) (records []model.Record, err error) {
	var content string
	if text, err := io.ReadAll(file); err == nil {
		content = string(text)
	}
	t := model.Record{}
	reg := regexp.MustCompile(`(?m)^requirepass\s+(.*)`)
	result := reg.FindAllStringSubmatch(content, -1)
	for _, passwd := range result {
		t.Username = "None"
		// 这里匹配到的结果必然是(requirepass password) 形式,
		// 因此Split(passwd[0], " ")[1] 用于获取密码
		t.Password = passwd[1]
		records = append(records, t)
	}
	return records, nil
}

func init() {
	mod := &redisService{}
	ServiceMatcherMap["redis"] = "plain"
	mod.name = "redis"
	mod.filepath = []string{"/etc/redis/redis.conf", "/etc/redisc.conf"}
	Register(mod)
}
