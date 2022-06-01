package plain

import (
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/hash"
)

type Plain struct {
	plain string
}

func (i *Plain) ID() string {
	return "Plain"
}

func (i *Plain) Plain() (plain string) {
	return i.plain
}

func (i *Plain) Match(dict []string) (guess string, err error) {
	for _, item := range dict {
		if item == i.plain {
			return item, nil
		}
	}
	return "", hash.ErrNotMatch
}
func New(plain string) (hash.Hash, error) {
	return &Plain{plain: plain}, nil
}
