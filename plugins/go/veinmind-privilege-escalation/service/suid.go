package service

import (
	"fmt"
	api "github.com/chaitin/libveinmind/go"
	"os"
	"strings"
	"syscall"
)

func SuidCheck(fs api.FileSystem, content os.FileInfo, filename string) (bool, error) {
	return isBelongToRoot(content) && isContainSUID(content), nil
}

func isBelongToRoot(content os.FileInfo) bool {
	uid := content.Sys().(*syscall.Stat_t).Uid
	if uid == uint32(0) {
		return true
	}
	return false
}

func isContainSUID(content os.FileInfo) bool {
	res := fmt.Sprintf("%o", uint32(content.Mode()))
	if strings.HasPrefix(res, "40000") {
		return true
	}
	return false
}
func init() {
	ImageCheckFuncMap["suid"] = SuidCheck
	ContainerCheckFuncMap["suid"] = SuidCheck
}
