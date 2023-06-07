package service

import (
	api "github.com/chaitin/libveinmind/go"
	"os"
)

func LimitedSuidCheck(fs api.FileSystem, content os.FileInfo, filename string) (bool, error) {
	return isBelongToRoot(content) && isContainSUID(content), nil
}

func init() {
	ImageCheckFuncMap["limited-suid"] = LimitedSuidCheck
	ContainerCheckFuncMap["limited-suid"] = LimitedSuidCheck
}
