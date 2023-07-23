package analyzer

import (
	"io/fs"
	"syscall"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/passwd"
	"github.com/chaitin/veinmind-common-go/service/report/event"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-trace/pkg/security"
)

func init() {
	Group = append(Group, &FileAnalyzer{})
}

type FileAnalyzer struct {
	event []*event.TraceEvent
}

func (fa *FileAnalyzer) Scan(container api.Container) {
	fa.event = make([]*event.TraceEvent, 0)
	fa.scanFilePerm(container)
	fa.scanTrace(container)
	fa.scanUser(container)
}

func (fa *FileAnalyzer) scanFilePerm(container api.Container) {
	// 检测敏感文件权限配置
	for dir, perm := range security.SensitiveDirPerm {
		if info, err := container.Stat(dir); err == nil {
			// check uid first
			sys := info.Sys()
			if stat, ok := sys.(*syscall.Stat_t); ok && stat.Uid != perm.Uid {
				fa.event = append(fa.event, &event.TraceEvent{
					Name:        "Perm of Sensitive File",
					From:        "File",
					Path:        dir,
					Description: "Incorrect user permissions for sensitive files",
					Detail:      string(stat.Uid),
					Level:       event.Medium,
				})
			}
			// check perm next
			if perm.Mode != 0 && info.Mode() != perm.Mode {
				fa.event = append(fa.event, &event.TraceEvent{
					Name:        "Perm of Sensitive File",
					From:        "File",
					Path:        dir,
					Description: "Incorrect file permissions for sensitive files",
					Detail:      info.Mode().String() + " (should " + perm.Mode.String() + ")",
					Level:       event.Medium,
				})
			}
		}
	}
}

func (fa *FileAnalyzer) scanTrace(container api.Container) {
	// 检测文件痕迹
	container.Walk("/tmp", func(path string, info fs.FileInfo, err error) error {
		for _, re := range security.CDKTrace {
			if re.MatchString(path) {
				fa.event = append(fa.event, &event.TraceEvent{
					Name:        "Attack Trace",
					From:        "File",
					Path:        path,
					Description: "Found traces of intrusion, it is highly likely that the container has been invaded",
					Detail:      "cdk trace file",
					Level:       event.High,
				})
			}
		}
		return nil
	})
	container.Walk("/mnt", func(path string, info fs.FileInfo, err error) error {
		for _, re := range security.CDKTrace {
			if re.MatchString(path) {
				fa.event = append(fa.event, &event.TraceEvent{
					Name:        "Attack Trace",
					From:        "File",
					Path:        path,
					Description: "Found traces of intrusion, it is highly likely that the container has been invaded",
					Detail:      "cdk trace file",
					Level:       event.High,
				})
			}
		}
		return nil
	})
}

func (fa *FileAnalyzer) scanUser(container api.Container) {
	// 检查用户passwd
	entries, err := passwd.ParseFilesystemPasswd(container)
	entryMap := make(map[string]string, 0)
	if err != nil {
		return
	}
	for _, e := range entries {
		// 1. check uid=0 but not root user
		if e.Uid == "0" && e.Username != "root" {
			fa.event = append(fa.event, &event.TraceEvent{
				Name:        "Abnormal user",
				From:        "File",
				Path:        "/etc/passwd",
				Description: "Abnormal user detected which uid=0 but not root",
				Detail:      e.Username + ";uid:" + e.Uid + ";gid:" + e.Gid,
				Level:       event.Medium,
			})
		}
		// 2. check gid=0 but not root user
		//if e.Gid == "0" && e.Username != "root" {
		//	fa.event = append(fa.event, &event.TraceEvent{
		//		Name:        "Abnormal user",
		//		From:        "File",
		//		Path:        "/etc/passwd",
		//		Description: "Abnormal user detected which gid=0 but not root",
		//		Detail:      e.Username + ";uid:" + e.Uid + ";gid:" + e.Gid,
		//		Level:       event.Medium,
		//	})
		//}
		// 3. check same uid user
		if _, ok := entryMap[e.Uid]; ok && e.Username != entryMap[e.Uid] {
			fa.event = append(fa.event, &event.TraceEvent{
				Name:        "Abnormal user",
				From:        "File",
				Path:        "/etc/passwd",
				Description: "Abnormal user detected which same uid",
				Detail:      e.Username + ";uid:" + e.Uid + ";gid:" + e.Gid,
				Level:       event.Medium,
			})
		} else {
			entryMap[e.Uid] = e.Username
		}

	}
}

func (fa *FileAnalyzer) Result() []*event.TraceEvent {
	return fa.event
}
