package service

import (
	"io/fs"
	"strings"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report/event"
)

// sshdBackdoorCheck 通过扫描文件软连接实现 sshd 软连接后门检测
func sshdBackdoorCheck(apiFileSystem api.FileSystem) (bool, []*event.BackdoorDetail) {
	var rootokList = []string{"su", "chsh", "chfn", "runuser"}
	check := false
	var res []*event.BackdoorDetail
	apiFileSystem.Walk("/", func(path string, info fs.FileInfo, err error) error {
		lstat, err := apiFileSystem.Lstat(path)
		if err != nil {
			return err
		}

		// 检查文件的软连接
		if lstat.Mode()&fs.ModeSymlink == fs.ModeSymlink {
			fLink, err := apiFileSystem.Readlink(path)
			if err != nil {
				return err
			}
			fExeName := path[strings.LastIndex(path, "/")+1:]
			fLinkExeName := fLink[strings.LastIndex(fLink, "/")+1:]
			if ContainsString(rootokList, fExeName) && fLinkExeName == "sshd" {
				check = true
				fileDetail, err := file2FileDetail(info, path)
				if err != nil {
					return nil
				}
				res = append(res, &event.BackdoorDetail{
					FileDetail:  fileDetail,
					Content:     fLink,
					Description: "sshd backdoor",
				})
			}
		}
		return nil
	})

	return check, res
}

func ContainsString(array []string, str string) bool {
	for _, s := range array {
		if s == str {
			return true
		}
	}
	return false
}

func init() {
	ImageCheckFuncMap["sshd"] = sshdBackdoorCheck
	ContainerCheckFuncMap["sshd"] = sshdBackdoorCheck
}
