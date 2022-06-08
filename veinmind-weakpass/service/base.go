package service

import (
	"fmt"
	"io"

	"github.com/chaitin/veinmind-tools/veinmind-weakpass/dict"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/hash"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/model"
)

// 对于每一个服务,需要对应一个爆破方法
// 服务需要在init函数指定与之对应的hash算法
var ServiceMatcherMap = make(map[string]string)

type IService interface {
	Name() string
	FilePath() []string
	GetRecords(file io.Reader) (records []model.Record, err error)
}

func GetDict(service string) (results []string) {
	results = append(results, dict.DictMap[service]...)
	results = append(results, dict.DictMap["base"]...)
	return results
}

func GetHash(service string) (hashI hash.Hash, err error) {
	for _, item := range hash.All {
		if item.ID() == ServiceMatcherMap[service] {
			hashI = item
			break
		}
	}
	if hashI == nil {
		return nil, fmt.Errorf("hash for service %s not found", service)
	}
	return hashI, nil
}
