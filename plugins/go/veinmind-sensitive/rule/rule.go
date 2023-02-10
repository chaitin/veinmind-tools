package rule

import (
	"errors"
	"io"

	"github.com/BurntSushi/toml"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/conf"
	"github.com/gogf/gf/errors/gerror"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-sensitive/embed"
)

type Config struct {
	WhiteList WhiteList `json:"white_list" toml:"white_list"`
	Rule      []Rule    `json:"rule" toml:"rule"`
	MIMEMap   map[string]bool
}

type WhiteList struct {
	PathPattern []string `json:"path_pattern" toml:"path_pattern"`
}

type Rule struct {
	Id              int64  `json:"id" toml:"id"`
	Name            string `json:"name" toml:"name"`
	Description     string `json:"description" toml:"description"`
	FilePathPattern string `json:"file_path_pattern" toml:"file_path_pattern"`
	Level           string `json:"level" toml:"level"`
	MIME            string `json:"mime" toml:"mime"`
	MatchPattern    string `json:"match_pattern" toml:"match_pattern"`
	Env             string `json:"env" toml:"env"`
}

var config *Config

func loadConfigFromService() (*Config, error) {
	data, err := conf.DefaultConfClient().Pull(conf.Sensitive)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = toml.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func loadConfigFromEmbed() (*Config, error) {
	fp, err := embed.FS.Open("rules.toml")
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(fp)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = toml.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func SingletonConf() *Config {
	return config
}

func Init() {
	conf, err := loadConfigFromService()
	if err != nil {
		c, e := loadConfigFromEmbed()
		if e != nil {
			e = gerror.Wrap(e, err.Error())
			log.Errorf("%+v", e)
			return
		}

		conf = c
	}
	config = conf

	if conf == nil {
		panic(errors.New("rule: can't init sensitiveConfig from service or local"))
	}

	if conf.MIMEMap == nil {
		conf.MIMEMap = make(map[string]bool)
	}

	for _, r := range conf.Rule {
		if r.MIME != "" {
			conf.MIMEMap[r.MIME] = true
		}
	}
}
