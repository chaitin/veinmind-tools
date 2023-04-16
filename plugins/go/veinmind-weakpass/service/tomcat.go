package service

import (
	"errors"
	"io"

	"github.com/beevik/etree"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/model"
)

type tomcatService struct {
	name     string
	filepath []string
}

func (i *tomcatService) Name() string {
	return i.name
}

func (i *tomcatService) FilePath() (paths []string) {
	return i.filepath
}
func (i *tomcatService) GetRecords(file io.Reader) (records []model.Record, err error) {
	doc := etree.NewDocument()
	if _, err := doc.ReadFrom(file); err != nil {
		return records, err
	}
	root := doc.SelectElement("tomcat-users")
	if root == nil {
		return records, errors.New("config file format error")
	}
	token := root.FindElements("user")
	if token == nil {
		return records, errors.New("config file format error")
	}
	t := model.Record{}
	for _, res := range token {
		t.Username = res.SelectAttr("username").Value
		t.Password = res.SelectAttr("password").Value
		t.Attributes = map[string]string{"roles": res.SelectAttr("roles").Value}
		records = append(records, t)
	}
	return records, nil
}

func init() {
	mod := &tomcatService{}
	ServiceMatcherMap["tomcat"] = "plain"
	mod.name = "tomcat"
	mod.filepath = []string{"/usr/local/tomcat/conf/tomcat-users.xml", "/etc/redis/tomcat-users.xml"}
	Register("tomcat", mod)
}
