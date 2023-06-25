package service

import (
	"io"
	"io/fs"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report/event"
)

func cronBackdoorCheck(apiFileSystem api.FileSystem) (bool, []*event.BackdoorDetail) {
	cronDirList := []string{"/var/spool/cron/", "/etc/cron.d/", "/etc/cron.daily", "/etc/cron.weekly/", "/etc/cron.hourly", "/etc/cron.monthly"}
	check := false
	var res []*event.BackdoorDetail
	for _, cronDir := range cronDirList {
		apiFileSystem.Walk(cronDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
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
					Description: "cronjob backdoor",
				})
			}
			return nil
		})
	}
	return check, res
}

func init() {
	ImageCheckFuncMap["cron"] = cronBackdoorCheck
	ContainerCheckFuncMap["cron"] = cronBackdoorCheck
}
