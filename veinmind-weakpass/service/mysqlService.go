package service

import (
	"io"

	"github.com/chaitin/veinmind-tools/veinmind-weakpass/extractor"
)

type mysqlService struct {
	name     string
	filepath []string
	extractor.Extractor
}

func (i *mysqlService) Name() string {
	return "mysql"
}

func (i *mysqlService) FilePath() (paths []string) {
	return i.filepath
}
func (i *mysqlService) GetRecords(file io.Reader) (records []extractor.Record, err error) {
	Extractor, err := GetExtractor("mysql")
	if err != nil {
		return records, err
	}
	// 从文件中获取密码相关的记录
	records, err = Extractor.Extract(file)
	if err != nil {
		return records, err
	}
	return records, nil
}

func init() {
	mod := &mysqlService{}
	mod.name = "mysql"
	mod.filepath = []string{"/var/lib/mysql/mysql.ibd", "/var/lib/mysql/mysql2.ibd"}
	Register(mod)
}
