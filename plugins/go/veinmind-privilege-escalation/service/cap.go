package service

import "C"

import (
	"os"

	api "github.com/chaitin/libveinmind/go"
)

// CapCheck 检测二进制文件是否有`CAP_SETUID`权限
func CapCheck(fs api.FileSystem, content os.FileInfo, filename string) (bool, error) {
	// TODO
	return false, nil
}

func init() {
	ImageCheckFuncMap["capabilities"] = CapCheck
	ContainerCheckFuncMap["capabilities"] = CapCheck
}
