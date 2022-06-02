package redis

import (
	"io"
	"regexp"
	"strings"

	"github.com/chaitin/veinmind-tools/veinmind-weakpass/extractor"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/hash/plain"
)

type Redis struct {
}

func (i *Redis) Meta() extractor.Meta {
	return extractor.Meta{Service: "redis"}
}
func (i *Redis) Extract(file io.Reader) (records []extractor.Record, err error) {
	var content string
	if text, err := io.ReadAll(file); err == nil {
		content = string(text)
	}
	t := extractor.Record{}
	reg := regexp.MustCompile(`[^# |#]requirepass .*`)
	result := reg.FindAllStringSubmatch(content, -1)
	for _, passwd := range result {
		t.Username = "None"
		h, _ := plain.New(strings.Split(passwd[0], " ")[1])
		t.Password = h
		records = append(records, t)
	}
	return records, nil
}
