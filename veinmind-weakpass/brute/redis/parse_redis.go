package redis

import (
	"io"
	"regexp"
	"strings"
)

type Redis struct {
	Filepath string
	Password string
}

func ParseRedisFile(redisFile io.Reader) (rediss []Redis, err error) {
	var content string
	if text, err := io.ReadAll(redisFile); err == nil {
		content = string(text)
	}
	redis := Redis{}
	reg := regexp.MustCompile(`[^# |#]requirepass .*`)
	result := reg.FindAllStringSubmatch(content, -1)
	for _, passwd := range result {
		redis.Password = strings.Split(passwd[0], " ")[1]
		rediss = append(rediss, redis)
	}

	return rediss, nil
}
