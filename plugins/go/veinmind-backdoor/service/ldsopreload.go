package service

import (
	"io"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report/event"
)

func ldsopreloadBackdoorCheck(apiFileSystem api.FileSystem) (bool, []*event.BackdoorDetail) {
	filePath := "/etc/ld.so.preload"
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
		fileDetail, err := file2FileDetail(fileInfo, filePath)
		if err != nil {
			return false, nil
		}
		check = true
		res = append(res, &event.BackdoorDetail{
			FileDetail:  fileDetail,
			Content:     content,
			Description: "ldsopreload conf backdoor",
		})
	}
	return check, res
}

func init() {
	ImageCheckFuncMap["ldsopreload"] = ldsopreloadBackdoorCheck
	ContainerCheckFuncMap["ldsopreload"] = ldsopreloadBackdoorCheck
}
