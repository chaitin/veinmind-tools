package service

import (
	"errors"
	"fmt"
)

var modules = make(map[string][]IService)

func Register(key string, p IService) {
	if p == nil {
		panic("Register service is nil")
	}
	modules[key] = append(modules[key], p)
}

// GetModuleByName 根据模块名获取对应Service列表
func GetModuleByName(modName string) ([]IService, error) {
	m, f := modules[modName]
	if f {
		return m, nil
	}
	return nil, errors.New(fmt.Sprintf("no mod named %s", modName))
}
