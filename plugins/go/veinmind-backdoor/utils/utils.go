package utils

import (
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report/event"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-backdoor/service"
)

func ImagesScanRun(fs api.Image) []*event.BackdoorDetail {
	var result = make([]*event.BackdoorDetail, 0)
	//调用对应模块检测是否有后门
	for _, checkFunc := range service.ImageCheckFuncMap {
		risk, detail := checkFunc(fs)
		// 如果有风险，将风险信息添加到result中
		if risk {
			result = append(result, detail...)
		}
	}
	return result
}

func ContainersScanRun(fs api.Container) []*event.BackdoorDetail {
	var result = make([]*event.BackdoorDetail, 0)
	//调用对应模块检测是否有后门
	for _, checkFunc := range service.ContainerCheckFuncMap {
		risk, detail := checkFunc(fs)
		// 如果有风险，将风险信息添加到result中
		if risk {
			result = append(result, detail...)
		}
	}
	return result
}
