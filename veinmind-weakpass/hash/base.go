package hash

import (
	"errors"
)

type Hash interface {
	// 密码匹配方式
	Match(dict []string) (plain string, err error)
	// 获取密码明文
	Plain() (plain string)
	// 加密算法的ID
	ID() string
}

var ErrNotMatch = errors.New("password and hash does not match")
var ErrCanNotGetPlain = errors.New("this hash can not get plain")
