package dict

import (
	"bufio"

	"github.com/chaitin/veinmind-tools/veinmind-weakpass/dict/embed"
)

var DictMap = make(map[string][]string, 10)

func Newdict(path string) (passDict []string) {
	passDictFile, err := embed.EmbedFS.Open(path)
	defer passDictFile.Close()
	if err != nil {
		panic(path + " create failed")
	}
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
}
