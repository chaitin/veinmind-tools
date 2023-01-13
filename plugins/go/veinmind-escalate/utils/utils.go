package utils

import (
	_ "encoding/json"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-escalate/pkg"
	_ "github.com/docker/docker/api/types/mount"
)

type checkFunc func(api.FileSystem) error

var (
	ImageCheckList = []checkFunc{
		pkg.UnsafeSuidCheck,
		pkg.CheckEmptyPasswdRoot,
		pkg.SudoFileCheck,
		pkg.UnsafePrivCheck,
	}
	ContainerCheckList = []checkFunc{
		pkg.ContainerUnsafeMount,
		pkg.ContainerUnsafeCapCheck,
		pkg.ContainerCVECheck,
		pkg.ContainerDockerAPiCheck,

		pkg.UnsafeSuidCheck,
		pkg.CheckEmptyPasswdRoot,
		pkg.SudoFileCheck,
		pkg.UnsafePrivCheck,
	}
)

func ImagesScanRun(fs api.Image) error {
	for _, opt := range ImageCheckList {
		opt(fs)
	}
	return pkg.GenerateImageRoport(fs)

}

func ContainersScanRun(fs api.Container) error {
	for _, opt := range ContainerCheckList {
		opt(fs)
	}
	return pkg.GenerateContainerRoport(fs)
}
