package tomcat

import (
	"errors"
	"io"

	"github.com/beevik/etree"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/extractor"

	"github.com/chaitin/veinmind-tools/veinmind-weakpass/hash/plain"
)

type Tomcat struct {
}

func (i *Tomcat) Meta() extractor.Meta {
	return extractor.Meta{Service: "tomcat"}
}
func (i *Tomcat) Extract(file io.Reader) (records []extractor.Record, err error) {
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
	t := extractor.Record{}
	for _, res := range token {
		t.Username = res.SelectAttr("username").Value
		h, _ := plain.New(res.SelectAttr("password").Value)
		t.Password = h
		t.Attributes = map[string]string{"roles": res.SelectAttr("roles").Value}
		records = append(records, t)
	}
	return records, nil
}
