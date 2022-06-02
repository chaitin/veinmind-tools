package extractor

import (
	"io"

	"github.com/chaitin/veinmind-tools/veinmind-weakpass/hash"
)

type Meta struct {
	Service string
}

type Record struct {
	Username   string
	Password   hash.Hash
	Attributes map[string]string
}
type Extractor interface {
	Meta() Meta
	Extract(file io.Reader) ([]Record, error)
}
