package service

import (
	"io"
	"strings"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/pkg/innodb"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/model"
)

type mysqlService struct {
	name     string
	filepath []string
}

func (i *mysqlService) Name() string {
	return i.name
}

func (i *mysqlService) FilePath() (paths []string) {
	return i.filepath
}
func (i *mysqlService) GetRecords(file io.Reader) (records []model.Record, err error) {
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
		tmp.Username = info.Name + "@" + info.Host
		tmp.Password = strings.ToLower(info.Password)
		records = append(records, tmp)
	}
	return records, nil
}

func init() {
	mod := &mysqlService{}
	ServiceMatcherMap["mysql"] = "mysql_native_password"
	mod.name = "mysql"
	mod.filepath = []string{"/var/lib/mysql/mysql.ibd", "/var/lib/mysql/mysql2.ibd"}
	Register(mod)
}
