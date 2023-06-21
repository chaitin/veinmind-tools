package service

import (
	"io"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report/event"
)

func aliasBackdoorCheck(fs api.FileSystem) (bool, *event.BackdoorDetail) {
	files := []string{"/root/.bashrc", "/root/.bash_profile", "/etc/bashrc", "/etc/profile"}

	for _, filename := range files {
		//校验文件是否存在
		fileInfo, err := fs.Stat(filename)
		if err != nil {
			return false, nil
		}
		//校验文件是否存在后门 )
		file, err := fs.Open(filename)
		if err != nil {
			return false, nil
		}
		defer file.Close()
		contents, err := io.ReadAll(file)
		risk, content := analysisStrings(string(contents))
		if risk {
			fileDetail, err := file2FileDetail(fileInfo, filename)
			if err != nil {
				return false, nil
			}
			return true, &event.BackdoorDetail{
				FileDetail:  fileDetail,
				Content:     content,
				Description: "alias backdoor",
			}
		}
	}
	return false, nil
}

func init() {
	ImageCheckFuncMap["alias"] = aliasBackdoorCheck
	ContainerCheckFuncMap["alias"] = aliasBackdoorCheck
}
