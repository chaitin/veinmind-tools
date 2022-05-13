package tomcat_passwd

import (
	"github.com/beevik/etree"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"io"
)

type Tomcat struct {
	Filepath string
	Username string
	Password string
	Role     string
}

// 需要考虑用户是否提供tomcat的目录,如果不提供需要自己找
func ParseTomcatFile(tomcatFile io.Reader) (tomcats []Tomcat, err error) {
	doc := etree.NewDocument()
	if _, err := doc.ReadFrom(tomcatFile); err != nil {
		log.Error(err)
	}
	root := doc.SelectElement("tomcat-users")
	if root == nil {
		log.Error("config file formate error")
	}
	token := root.FindElements("user")
	if token == nil {
		log.Error("config file formate error")
	}
	t := Tomcat{}
	for _, res := range token {
		t.Username = res.SelectAttr("username").Value
		t.Password = res.SelectAttr("password").Value
		t.Role = res.SelectAttr("roles").Value
		tomcats = append(tomcats, t)
	}

	return tomcats, nil
}
