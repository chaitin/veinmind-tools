package service

import (
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report/event"
)

func sshdBackdoorCheck(fs api.FileSystem) (bool, *event.BackdoorDetail) {
	// TODO
	return false, nil
}

func init() {
	ImageCheckFuncMap["sshd"] = sshdBackdoorCheck
	ContainerCheckFuncMap["sshd"] = sshdBackdoorCheck
}
