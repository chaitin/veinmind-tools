package service

import (
	"io"
	"strings"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/model"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/pkg/innodb"
)

var _ IService = (*mysql8Service)(nil)

type mysql8Service struct {
}

func (i *mysql8Service) Name() string {
	return "mysql8"
}

func (i *mysql8Service) FilePath() (paths []string) {
	return []string{"/var/lib/mysql/mysql.ibd", "/var/lib/mysql/mysql2.ibd"}
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
		if strings.Contains(info.Password, innodb.EmptyPasswordPlaceholder) {
			continue
		}
		if info.Plugin != innodb.PluginNameNative {
			tmp.Password = info.Password
		} else {
			tmp.Password = strings.ToLower(info.Password)
		}
		tmp.Username = info.Name
		records = append(records, tmp)
	}
	return records, nil
}

func init() {
	mod := &mysql8Service{}
	ServiceMatcherMap[mod.Name()] = "mysql"
	Register("mysql", mod)
}
