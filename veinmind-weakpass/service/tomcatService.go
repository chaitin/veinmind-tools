package service

import (
	"io"

	"github.com/chaitin/veinmind-tools/veinmind-weakpass/extractor"
)

type tomcatService struct {
	name     string
	filepath []string
	extractor.Extractor
}

func (i *tomcatService) Name() string {
	return "tomcat"
}

func (i *tomcatService) FilePath() (paths []string) {
	return i.filepath
}
func (i *tomcatService) GetRecords(file io.Reader) (records []extractor.Record, err error) {
	Extractor, err := GetExtractor("tomcat")
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
	mod := &tomcatService{}
	mod.name = "tomcat"
	mod.filepath = []string{"/usr/local/tomcat/conf/tomcat-users.xml", "/etc/redis/tomcat-users.xml"}
	Register(mod)
}
