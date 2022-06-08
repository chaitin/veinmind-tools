package service

import (
	"errors"
	"fmt"
)

var modules = make(map[string]IService)

func Register(p IService) {
	if p == nil {
		panic("Register service is nil")
	}
	name := p.Name()
	if _, dup := modules[name]; dup {
		panic(fmt.Sprintf("Register called twice for service %s", name))
	}
	modules[name] = p
}

// GetModules 获取modules列表
func GetAllModules() map[string]IService {
	return modules
}

// GetModulesByName 根据模块名获取modules列表
func GetModuleByName(modName string) (IService, error) {
	m, f := modules[modName]
	if f {
		return m, nil
	}
	return nil, errors.New(fmt.Sprintf("no mod named %s", modName))
}
