package service

import (
	"io"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report/event"
)

func inetdBackdoorCheck(apiFileSystem api.FileSystem) (bool, []*event.BackdoorDetail) {
	filePath := "/etc/inetd.conf"
	check := false
	var res []*event.BackdoorDetail

	fileInfo, err := apiFileSystem.Stat(filePath)
	if err != nil {
		return false, nil
	}
	file, err := apiFileSystem.Open(filePath)
	if err != nil {
		return false, nil
	}
	contents, err := io.ReadAll(file)
	risk, content := analysisStrings(string(contents))
	if risk {
		check = true
		fileDetail, err := file2FileDetail(fileInfo, filePath)
		if err != nil {
			return false, nil
		}
		res = append(res, &event.BackdoorDetail{
			FileDetail:  fileDetail,
			Content:     content,
			Description: "inetd conf backdoor",
		})
	}
	return check, res
}

func init() {
	ImageCheckFuncMap["inetd"] = inetdBackdoorCheck
	ContainerCheckFuncMap["inetd"] = inetdBackdoorCheck
}
