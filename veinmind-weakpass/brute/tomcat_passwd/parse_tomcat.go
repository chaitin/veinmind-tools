package tomcat_passwd

import (
	"io"
	"regexp"
)
type Tomcat struct{
	Filepath string
	Username string
	Password string
	Role string
}
// 需要考虑用户是否提供tomcat的目录,如果不提供需要自己找
func ParseTomcatFile(tomcatFile io.Reader) (tomcats []Tomcat, err error) {
	var content string
	if text, err := io.ReadAll(tomcatFile); err == nil {
		content = string(text)
	}
	t := Tomcat{} 
	// user 标签下是否有这三个属性,如果不全要按行进行匹配,每次匹配一个属性即可
	reg := regexp.MustCompile(`<user username="(?s:(.*?))" password="(?s:(.*?))" roles="(?s:(.*?))"/>`)
	result := reg.FindAllStringSubmatch(content, -1)
	for _, text := range result {
		t.Username = text[1]
		t.Password = text[2]
		t.Role = text[3]
		tomcats = append(tomcats,t)
    }
	return tomcats, nil
}