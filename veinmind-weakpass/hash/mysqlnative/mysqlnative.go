package mysqlnative

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"strings"

	"github.com/chaitin/veinmind-tools/veinmind-weakpass/hash"
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
func (i *MysqlNative) Match(dict []string) (plain string, err error) {
	if strings.Contains(i.hash, "*") {
		for _, guess := range dict {
			r := sha1.Sum([]byte(guess))
			r = sha1.Sum(r[:])
			s := fmt.Sprintf("%x", r)
			if strings.Contains(i.hash, s) {
				return guess, nil
			}

		}
	}
	return plain, errors.New("invalid mysql_native_password format")
}

func New(rawHash string) (hash.Hash, error) {
	return &MysqlNative{hash: rawHash}, nil
}
