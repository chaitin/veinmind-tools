package service

import (
	"bufio"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/model"
	"io"
)

type FtpService struct {
	name     string
	filepath []string
}

func (i *FtpService) Name() string {
	return i.name
}

func (i *FtpService) FilePath() (paths []string) {
	return i.filepath
}
func (i *FtpService) GetRecords(file io.Reader) (records []model.Record, err error) {
	scanner := bufio.NewScanner(file)
	//文件奇数行为username 偶数行为password
	lineNum := 0
	s := model.Record{}
	for scanner.Scan() {
		userinfo := scanner.Text()
		if lineNum%2 == 0 {
			s.Username = userinfo
		} else {
			s.Password = userinfo
			records = append(records, s)
			s = model.Record{}
		}
		lineNum++
	}

	return records, nil
}

func init() {
	mod := &FtpService{}
	ServiceMatcherMap["ftp"] = "plain"
	mod.name = "ftp"
	mod.filepath = []string{"/etc/vsftpd/virtual_users.txt"}
	Register(mod)
}
