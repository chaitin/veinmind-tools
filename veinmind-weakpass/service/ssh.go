package service

import (
	"bufio"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"io"
	"strings"

	"github.com/chaitin/veinmind-tools/veinmind-weakpass/model"
)

type sshService struct {
	name     string
	filepath []string
}

func (i *sshService) Name() string {
	return i.name
}

func (i *sshService) FilePath() (paths []string) {
	return i.filepath
}
func (i *sshService) GetRecords(file io.Reader) (records []model.Record, err error) {
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
	mod := &sshService{
		name:     "ssh",
		filepath: []string{"/etc/shadow"},
	}
	ServiceMatcherMap["ssh"] = "shadow"
	Register(mod)
}
