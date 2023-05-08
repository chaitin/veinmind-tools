package service

import (
	"io"
	"strings"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/pkg/myisam"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/model"
)

type mysql5Service struct {
	name     string
	filepath []string
}

func (i *mysql5Service) Name() string {
	return i.name
}

func (i *mysql5Service) FilePath() (paths []string) {
	return i.filepath
}
func (i *mysql5Service) GetRecords(file io.Reader) (records []model.Record, err error) {
	mysqlInfos, err := myisam.ParseUserFile(file)
	if err != nil {
		return records, err
	}
	tmp := model.Record{}
	for _, info := range mysqlInfos {
		tmp.Username = info.Name
		tmp.Password = strings.ToLower(info.Password)
		records = append(records, tmp)
	}
	return records, nil
}

func init() {
	mod := &mysql5Service{}
	ServiceMatcherMap["mysql5"] = "mysql_native_password"
	mod.name = "mysql5"
	mod.filepath = []string{"/var/lib/mysql/mysql/user.MYD"}
	Register("mysql", mod)
}
