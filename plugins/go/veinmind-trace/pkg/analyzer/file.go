package analyzer

import (
	"syscall"

	api "github.com/chaitin/libveinmind/go"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-trace/pkg/security"
)

type FileAnalyzer struct {
	Object api.Container
	//Event    []
}

func (fa *FileAnalyzer) Scan() {
	fa.scanFilePerm()

}

func (fa *FileAnalyzer) scanFilePerm() {
	for dir, perm := range security.SensitiveDirPerm {
		if info, err := fa.Object.Stat(dir); err == nil {
			// check uid first
			sys := info.Sys()
			if stat, ok := sys.(*syscall.Stat_t); ok && stat.Uid != perm.Uid {
				// todo：add event
			}
			// check perm next
			if perm.Mode != 0 && info.Mode() != perm.Mode {
				// todo：add event
			}
		}
	}
}
