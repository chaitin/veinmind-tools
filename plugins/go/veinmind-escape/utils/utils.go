package utils

import (
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report/event"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-escape/pkg"
)

func ImagesScanRun(fs api.Image) []*event.EscapeDetail {
	var result = make([]*event.EscapeDetail, 0)
	for _, check := range pkg.ImageCheckList {
		res, err := check(fs)
		if err != nil {
			continue
		}
		result = append(result, res...)
	}
	return result
}

func ContainersScanRun(fs api.Container) []*event.EscapeDetail {
	var result = make([]*event.EscapeDetail, 0)
	for _, check := range pkg.ContainerCheckList {
		res, err := check(fs)
		if err != nil {
			continue
		}
		result = append(result, res...)
	}
	return result
}
