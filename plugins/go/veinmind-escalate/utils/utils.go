package utils

import (
	_ "encoding/json"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-escalate/pkg"
	_ "github.com/docker/docker/api/types/mount"
)

type checkFunc func(api.FileSystem) error

var (
	imageCheckList = []checkFunc{
		pkg.UnsafeSuidCheck,
		pkg.CheckEmptyPasswdRoot,
		pkg.SudoFileCheck,
		pkg.UnsafePrivCheck,
	}
	containerCheckList = []checkFunc{
		pkg.UnsafeSuidCheck,
		pkg.CheckEmptyPasswdRoot,
		pkg.UnsafeCapCheck,
		pkg.SudoFileCheck,
		pkg.UnsafePrivCheck,
		pkg.DetectContainerUnsafeMount,
	}
)

func ImagesScanRun(fs api.Image) error {
	for _, opt := range imageCheckList {
		opt(fs)
	}
	return pkg.GenerateImageRoport(fs)

}

func ContainersScanRun(fs api.Container) error {
	for _, opt := range containerCheckList {
		opt(fs)
	}
	return pkg.GenerateContainerRoport(fs)
}
