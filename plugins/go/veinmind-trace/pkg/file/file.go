package file

import (
	"syscall"

	api "github.com/chaitin/libveinmind/go"
)

func ScanFilePerm(container api.Container) {
	//1. do perm check
	for dir, perm := range sensitiveDirPerm {
		if info, err := container.Stat(dir); err == nil {
			// check uid first
			sys := info.Sys()
			if stat, ok := sys.(*syscall.Stat_t); ok && stat.Uid != perm.uid {
				// todo：add event
			}
			// check perm next
			if perm.mode != 0 && info.Mode() != perm.mode {
				// todo：add event
			}
		}
	}
	// 2.
}
