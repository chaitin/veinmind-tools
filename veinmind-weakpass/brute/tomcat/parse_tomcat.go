package tomcat

import (
	"errors"
	"github.com/beevik/etree"
	"io"
)

type Tomcat struct {
	Filepath string
	Username string
	Password string
	Role     string
}

func ParseTomcatFile(tomcatFile io.Reader) (tomcats []Tomcat, err error) {
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
	t := Tomcat{}
	for _, res := range token {
		t.Username = res.SelectAttr("username").Value
		t.Password = res.SelectAttr("password").Value
		t.Role = res.SelectAttr("roles").Value
		tomcats = append(tomcats, t)
	}

	return tomcats, nil
}
