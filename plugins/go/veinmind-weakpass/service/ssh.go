package service

import (
	"bufio"
	"io"
	"strings"

	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/model"
)

type SshService struct {
	name     string
	filepath []string
}

func (i *SshService) Name() string {
	return i.name
}

func (i *SshService) FilePath() (paths []string) {
	return i.filepath
}
func (i *SshService) GetRecords(file io.Reader) (records []model.Record, err error) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		userinfo := strings.Split(scanner.Text(), ":")
		if len(userinfo) != 9 {
			log.Error("service: shadow format error")
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
	mod := &SshService{}
	ServiceMatcherMap["ssh"] = "shadow"
	mod.name = "ssh"
	mod.filepath = []string{"/etc/shadow"}
	Register("ssh", mod)
}
