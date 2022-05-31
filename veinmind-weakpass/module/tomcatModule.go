package module

import (
	"errors"
	"github.com/beevik/etree"
	"io"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/dict"
)

type Tomcat struct {
	Module
}

func (this *Tomcat) Init(conf Config) (err error) {
	err = this.Module.Init(conf)
	this.specialDict = dict.Tomcatdict
	return err
}

func (this *Tomcat) Name() string {
	return this.name
}
func (this *Tomcat) ParsePasswdInfo(tomcatFile io.Reader) (tomcats []PasswdInfo, err error) {
	doc := etree.NewDocument()
	if _, err := doc.ReadFrom(tomcatFile); err != nil {
		return tomcats, err
	}
	root := doc.SelectElement("tomcat-users")
	if root == nil {
		return tomcats, errors.New("config file format error")
	}
	token := root.FindElements("user")
	if token == nil {
		return tomcats, errors.New("config file format error")
	}
	t := PasswdInfo{}
	for _, res := range token {
		t.Username = res.SelectAttr("username").Value
		t.Password = res.SelectAttr("password").Value
		tomcats = append(tomcats, t)
	}
	return tomcats, nil
}

func init() {
	mod := &Tomcat{}
	mod.name = "TOMCAT"
	mod.passwdType = TOMCAT
	mod.filePath = []string{"/usr/local/tomcat/conf/tomcat-users.xml"}
	Register(mod)
}
