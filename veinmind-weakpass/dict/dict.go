package dict

import (
	"bufio"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/embed"
)

var Passdict = Newdict("pass.dict")
var Tomcatdict = Newdict("tomcat.dict")
var Sshdict = Newdict("ssh.dict")
var Redisdict = Newdict("redis.dict")

func Newdict(path string) (passDict []string) {
	passDictFile, err := embed.EmbedFS.Open(path)
	if err != nil {
		panic(path+" create failed")
	}
	scanner := bufio.NewScanner(passDictFile)
	for scanner.Scan() {
		passDict = append(passDict, scanner.Text())
	}
	return passDict
}
