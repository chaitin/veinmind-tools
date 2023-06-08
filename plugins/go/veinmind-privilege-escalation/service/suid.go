package service

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	api "github.com/chaitin/libveinmind/go"
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
	// 4000 is SUID, 6000 is SUID and SGID
	if strings.HasPrefix(res, "4000") || strings.HasPrefix(res, "6000") {
		return true
	}
	return false
}
func init() {
	ImageCheckFuncMap["suid"] = SuidCheck
	ContainerCheckFuncMap["suid"] = SuidCheck
}
