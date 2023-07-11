package service

import (
	"io"
	"io/fs"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report/event"
)

// startupBackdoorCheck 检测启动项是否有后门
func startupBackdoorCheck(apiFileSystem api.FileSystem) (bool, []*event.BackdoorDetail) {
	startupDirs := []string{"/etc/init.d/", "/etc/rc.d/", "/etc/rc.local/", "/usr/local/etc/rc.d/", "/usr/local/etc/rc.local/", "/etc/conf.d/local.start/", "/etc/inittab/", "/etc/systemd/system/"}
	check := false
	var res []*event.BackdoorDetail
	for _, startupDir := range startupDirs {
		apiFileSystem.Walk(startupDir, func(path string, info fs.FileInfo, err error) error {
			file, err := apiFileSystem.Open(path)
			if err != nil {
				return nil
			}
			defer file.Close()
			contents, err := io.ReadAll(file)
			risk, content := analysisStrings(string(contents))
			if risk {
				check = true
				fileDetail, err := file2FileDetail(info, path)
				if err != nil {
					return nil
				}
				res = append(res, &event.BackdoorDetail{
					FileDetail:  fileDetail,
					Content:     content,
					Description: "startup backdoor",
				})
			}
			return nil
		})
	}

	return check, res
}

func init() {
	ImageCheckFuncMap["startup"] = startupBackdoorCheck
	ContainerCheckFuncMap["startup"] = startupBackdoorCheck
}
