package service

import (
	"io"
	"io/fs"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report/event"
)

func xinetdBackdoorCheck(apiFileSystem api.FileSystem) (bool, []*event.BackdoorDetail) {
	xinetdDir := "/etc/xinetd.conf/"
	check := false
	var res []*event.BackdoorDetail

	apiFileSystem.Walk(xinetdDir, func(path string, info fs.FileInfo, err error) error {
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
				Description: "xinetd backdoor",
			})
		}
		return nil
	})
	return check, res
}

func init() {
	ImageCheckFuncMap["xinetd"] = xinetdBackdoorCheck
	ContainerCheckFuncMap["xinetd"] = xinetdBackdoorCheck
}
