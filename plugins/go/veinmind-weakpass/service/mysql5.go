package service

import (
	"io"
	"strings"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/pkg/myisam"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/model"
)

type mysql5Service struct {
}

func (i *mysql5Service) Name() string {
	return "mysql5"
}

func (i *mysql5Service) FilePath() (paths []string) {
	return []string{"/var/lib/mysql/mysql/user.MYD"}
}
func (i *mysql5Service) GetRecords(file io.Reader) (records []model.Record, err error) {
	mysqlInfos, err := myisam.ParseUserFile(file)
	if err != nil {
		return records, err
	}
	tmp := model.Record{}
	for _, info := range mysqlInfos {
		if info.Password == myisam.EmptyPasswordPlaceholder || info.Host == myisam.LocalHost || info.Host == "127.0.0.1" {
			continue
		}
		tmp.Username = info.Name
		tmp.Password = strings.ToLower(info.Password)
		records = append(records, tmp)
	}
	return records, nil
}

func init() {
	mod := &mysql5Service{}
	ServiceMatcherMap[mod.Name()] = "mysql"
	Register("mysql", mod)
}
