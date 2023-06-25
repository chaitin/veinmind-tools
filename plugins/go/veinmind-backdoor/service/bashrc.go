package service

import (
	"io"
	"io/fs"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report/event"
)

func bashrcBackdoorCheck(apiFileSystem api.FileSystem) (bool, []*event.BackdoorDetail) {
	filePaths := []string{"/root/.bashrc", "root/.tcshrc", "/root/.bash_profile", "root/.cshrc", "/etc/.bashrc", "/etc/bashrc", "/etc/profile"}
	profileDir := "/etc/profile.d"
	homeDir := "/home"
	homeFiles := []string{".bashrc", ".bash_profile", ".tcshrc", ".cshrc"}

	check := false
	var res []*event.BackdoorDetail
	// 校验/root和/etc下的shell配置文件环境变量
	for _, filepath := range filePaths {
		//校验文件是否存在
		fileInfo, err := apiFileSystem.Stat(filepath)
		if err != nil {
			return false, nil
		}
		//校验文件是否存在后门
		file, err := apiFileSystem.Open(filepath)
		if err != nil {
			return false, nil
		}
		defer file.Close()
		contents, err := io.ReadAll(file)
		risk, content := analysisStrings(string(contents))
		if risk {
			check = true
			fileDetail, err := file2FileDetail(fileInfo, filepath)
			if err != nil {
				return false, nil
			}
			res = append(res, &event.BackdoorDetail{
				FileDetail:  fileDetail,
				Content:     content,
				Description: "env backdoor",
			})
		}
	}

	// 校验/etc/profile.d下的shell配置文件环境变量
	apiFileSystem.Walk(profileDir, func(path string, info fs.FileInfo, err error) error {
		//校验文件是否存在后门
		file, err := apiFileSystem.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()
		contents, err := io.ReadAll(file)
		risk, content := analysisStrings(string(contents))
		if risk {
			fileDetail, err := file2FileDetail(info, path)
			if err != nil {
				return nil
			}
			res = append(res, &event.BackdoorDetail{
				FileDetail:  fileDetail,
				Content:     content,
				Description: "env backdoor",
			})
		}
		return nil
	})

	// 校验/home下用户的shell配置文件环境变量
	apiFileSystem.Walk(homeDir, func(path string, info fs.FileInfo, err error) error {
		for _, filename := range homeFiles {
			if info.Name() == filename {
				//校验文件是否存在后门
				file, err := apiFileSystem.Open(path)
				if err != nil {
					return nil
				}
				defer file.Close()
				contents, err := io.ReadAll(file)
				risk, content := analysisStrings(string(contents))
				if risk {
					fileDetail, err := file2FileDetail(info, path)
					if err != nil {
						return nil
					}
					res = append(res, &event.BackdoorDetail{
						FileDetail:  fileDetail,
						Content:     content,
						Description: "env backdoor",
					})
				}
			}
		}
		return nil
	})

	return check, res
}

func init() {
	ImageCheckFuncMap["bashrc"] = bashrcBackdoorCheck
	ContainerCheckFuncMap["bashrc"] = bashrcBackdoorCheck
}
