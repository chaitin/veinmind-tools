package hash

import (
	"crypto/sha1"
	"fmt"
	"strings"
)

type MysqlNative struct {
	hash string
}

func (i *MysqlNative) ID() string {
	return "mysql_native_password"
}

func (i *MysqlNative) Plain() string {
	return string(i.hash)
}
func (i *MysqlNative) Match(hash, guess string) (result string, flag bool) {
	if strings.Contains(hash, "*") {
		r := sha1.Sum([]byte(guess))
		r = sha1.Sum(r[:])
		s := fmt.Sprintf("%x", r)
		if strings.Contains(hash, s) {
			return guess, true

		}
	}
	return result, false
}
