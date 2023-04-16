package service

import (
	"bufio"
	"io"
	"strings"

	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/model"
)

type proftpdService struct {
	name     string
	filepath []string
}

func (i *proftpdService) Name() string {
	return i.name
}

func (i *proftpdService) FilePath() (paths []string) {
	return i.filepath
}

// GetRecords TODO: 对同一个配置文件支持不同的Hash算法
func (i *proftpdService) GetRecords(file io.Reader) (records []model.Record, err error) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		userinfo := strings.Split(scanner.Text(), ":")
		if len(userinfo) != 7 {
			log.Error("service: proftpd format error")
			continue
		}

		s := model.Record{}
		s.Username = userinfo[0]
		s.Password = userinfo[1]
		records = append(records, s)
	}
	return records, nil
}

func init() {
	mod := &proftpdService{}
	ServiceMatcherMap["proftpd"] = "shadow"
	mod.name = "proftpd"
	mod.filepath = []string{"/etc/proftpd/passwd", "/etc/proftpd/ftppasswd", "/etc/proftpd/ftpd.passwd"}
	Register("ftp", mod)
}
