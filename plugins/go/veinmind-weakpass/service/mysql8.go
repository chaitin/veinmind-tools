package service

import (
	"io"
	"strings"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/pkg/innodb"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/model"
)

type mysql8Service struct {
	name     string
	filepath []string
}

func (i *mysql8Service) Name() string {
	return i.name
}

func (i *mysql8Service) FilePath() (paths []string) {
	return i.filepath
}
func (i *mysql8Service) GetRecords(file io.Reader) (records []model.Record, err error) {
	page, err := innodb.FindUserPage(file)
	if err != nil {
		return records, err
	}

	mysqlInfos, err := innodb.ParseUserPage(page.Pagedata)
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
	// TODO： Mysql8
	mod := &mysql8Service{}
	ServiceMatcherMap["mysql8"] = "mysql_native_password"
	mod.name = "mysql8"
	mod.filepath = []string{"/var/lib/mysql/mysql.ibd"}
	Register("mysql", mod)
}
