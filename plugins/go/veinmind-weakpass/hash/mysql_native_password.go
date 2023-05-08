package hash

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"strings"
)

type MysqlNative struct {
}

func (i *MysqlNative) ID() string {
	return "mysql_native_password"
}

func (i *MysqlNative) Match(hash, guess string) (flag bool, err error) {
	r := sha1.Sum([]byte(guess))
	r = sha1.Sum(r[:])
	s := fmt.Sprintf("%x", r)
	if strings.Contains(hash, s) {
		return true, nil
	}
	return false, errors.New("mysql_passwd: malformed entry ")
}
