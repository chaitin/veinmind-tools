package model

import (
	"io"
	"regexp"
	"strings"
)

type Redis struct {
	Module
}

func (this *Redis) Name() string {
	return this.name
}
func (this *Redis) Init(conf Config) error {
	this.Module.Init(conf)
	return nil
}
func (this *Redis) ParsePasswdInfo(redisFile io.Reader) (rediss []PasswdInfo, err error) {
	var content string
	if text, err := io.ReadAll(redisFile); err == nil {
		content = string(text)
	}
	redis := PasswdInfo{}
	reg := regexp.MustCompile(`[^# |#]requirepass .*`)
	result := reg.FindAllStringSubmatch(content, -1)
	for _, passwd := range result {
		redis.Password = strings.Split(passwd[0], " ")[1]
		rediss = append(rediss, redis)
	}
	return rediss, nil
}
func init() {
	mod := &Redis{}
	mod.passwdType = REDIS
	mod.name = "REDIS"
	mod.filePath = []string{"/etc/redis/redis.conf"}
	Register(mod)
}
