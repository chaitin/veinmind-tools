package dict

import (
	"bufio"
	"fmt"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/dict/embed"
)

var DictMap = make(map[string][]string, 10)

func Newdict(path string) (passDict []string) {
	passDictFile, err := embed.EmbedFS.Open(path)
	if err != nil {
		panic(fmt.Sprintf("create passdict failed cause %s", err))
	}
	defer passDictFile.Close()

	scanner := bufio.NewScanner(passDictFile)
	for scanner.Scan() {
		passDict = append(passDict, scanner.Text())
	}
	return passDict
}

func init() {
	DictMap["base"] = Newdict("pass.dict")
	DictMap["tomcat"] = Newdict("tomcat.dict")
	DictMap["ssh"] = Newdict("ssh.dict")
	DictMap["redis"] = Newdict("redis.dict")
	DictMap["ftp"] = Newdict("ftp.dict")
}
