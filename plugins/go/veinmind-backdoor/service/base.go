package service

import (
	"io/fs"
	"syscall"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report/event"
)

type CheckFunc func(fs api.FileSystem) (bool, []*event.BackdoorDetail)

var (
	ImageCheckFuncMap     = make(map[string]CheckFunc)
	ContainerCheckFuncMap = make(map[string]CheckFunc)
)

func file2FileDetail(info fs.FileInfo, path string) (event.FileDetail, error) {
	sys := info.Sys().(*syscall.Stat_t)

	return event.FileDetail{
		Path: path,
		Perm: info.Mode(),
		Size: info.Size(),
		Uid:  int64(sys.Uid),
		Gid:  int64(sys.Gid),
		Ctim: int64(sys.Ctim.Sec),
		Mtim: int64(sys.Mtim.Sec),
		Atim: int64(sys.Mtim.Sec),
	}, nil
}
