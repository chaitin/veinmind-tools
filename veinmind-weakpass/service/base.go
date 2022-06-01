package service

import (
	"fmt"
	"io"

	"github.com/chaitin/veinmind-tools/veinmind-weakpass/dict"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/extractor"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/extractor/all"
)

type IService interface {
	Name() string
	FilePath() []string
	GetRecords(file io.Reader) (records []extractor.Record, err error)
}

func GetExtractor(service string) (extractor extractor.Extractor, err error) {
	for _, e := range all.All {
		if e.Meta().Service == service {
			extractor = e
			break
		}
	}
	if extractor == nil {
		return nil, fmt.Errorf("extractor for service %s not found", service)
	}
	return extractor, nil
}

func GetDict(service string) (results []string) {
	results = append(results, dict.DictMap[service]...)
	results = append(results, dict.DictMap["base"]...)
	return results
}
