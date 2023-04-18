package service

import (
	"io"
	"regexp"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/model"
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
	} else {
		return records, err
	}
	t := model.Record{}
	reg := regexp.MustCompile(`(?m)^requirepass\s+[\"|\']?([^(\"|\'|\n)]+)[\"|\']?`)
	result := reg.FindAllStringSubmatch(content, -1)
	for _, passwd := range result {
		if len(passwd) != 2 {
			continue
		}
		t.Username = ""
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
	Register("redis", mod)
}
