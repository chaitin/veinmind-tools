package service

import (
	"io"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report/event"
)

func tcpWrapperBackdoorCheck(apiFileSystem api.FileSystem) (bool, []*event.BackdoorDetail) {
	filePaths := []string{"/etc/hosts.allow", "/etc/hosts.deny"}
	check := false
	var res []*event.BackdoorDetail
	for _, filePath := range filePaths {
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
				Description: "tcp wrapper backdoor",
			})
		}
	}

	return check, res
}

func init() {
	ImageCheckFuncMap["tcpWrapper"] = tcpWrapperBackdoorCheck
	ContainerCheckFuncMap["tcpWrapper"] = tcpWrapperBackdoorCheck
}
