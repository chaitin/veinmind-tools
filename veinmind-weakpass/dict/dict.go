package dict

import (
	"bufio"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/embed"
)

var Passdict = Newdict()

func Newdict() (passDict []string) {
	passDictFile, err := embed.EmbedFS.Open("pass.dict")
	if err != nil {
		panic("dict create failed")
	}
	scanner := bufio.NewScanner(passDictFile)
	for scanner.Scan() {
		passDict = append(passDict, scanner.Text())
	}
	return passDict
}
