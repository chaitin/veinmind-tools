package service

import (
	api "github.com/chaitin/libveinmind/go"
	"os"
)

type CheckFunc func(fs api.FileSystem, content os.FileInfo, filename string) (bool, error)

var (
	ImageCheckFuncMap     = make(map[string]CheckFunc)
	ContainerCheckFuncMap = make(map[string]CheckFunc)
)

const (
	SUDOREGEX string = "(\\w{1,})\\s\\w{1,}=\\(.*\\)\\s(.*)"
)
