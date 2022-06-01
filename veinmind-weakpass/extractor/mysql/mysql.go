package mysql

import (
	"io"
	"strings"

	"github.com/chaitin/veinmind-tools/veinmind-weakpass/hash/mysqlnative"

	"github.com/chaitin/veinmind-tools/veinmind-weakpass/extractor"
)

type Mysql struct {
}

func (i *Mysql) Meta() extractor.Meta {
	return extractor.Meta{Service: "mysql"}
}
func (this *Mysql) Extract(file io.Reader) (records []extractor.Record, err error) {
	page, err := FindUserPage(file)
	if err != nil {
		return records, err
	}

	mysqlInfos, err := ParseUserPage(page.Pagedata)
	if err != nil {
		return records, err
	}
	tmp := extractor.Record{}
	for _, info := range mysqlInfos {
		tmp.Username = info.Name + "@" + info.Host
		h, _ := mysqlnative.New(strings.ToLower(info.Password))
		tmp.Password = h
		records = append(records, tmp)
	}
	return records, nil
}
