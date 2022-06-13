package rule

import (
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-sensitive/embed"
	"github.com/chaitin/veinmind-tools/veinmind-common/go/service/conf"
	"github.com/gobwas/glob"
	"io/ioutil"
	"regexp"
	"strings"
)

type SensitiveConfig struct {
	WhiteList SensitiveWhiteList `json:"white_list" ,toml:"white_list"`
	Rules     []SensitiveRule    `json:"rules" ,toml:"rules"`
}

type SensitiveWhiteList struct {
	Paths     []string `json:"paths" ,toml:"paths"`
	PathsGlob []glob.Glob
}

type SensitiveRule struct {
	Id             int64  `json:"id" ,toml:"id"`
	Name           string `json:"name" ,toml:"name"`
	Description    string `json:"description" ,toml:"description"`
	Filepath       string `json:"filepath" ,toml:"filepath"`
	FilepathRegexp *regexp.Regexp
	Level          string `json:"level" ,toml:"level"`
	MIME           string `json:"mime" ,toml:"mime"`
	Match          string `json:"match" ,toml:"match"`
	MatchContains  string
	MatchRegex     *regexp.Regexp
	Env            string `json:"env" ,toml:"env"`
}

var sensitiveConfig *SensitiveConfig

func loadConfigFromService() (*SensitiveConfig, error) {
	confBytes, err := conf.DefaultConfClient().Pull(conf.Sensitive)
	if err != nil {
		return nil, err
	}

	confE := SensitiveConfig{}
	err = toml.Unmarshal(confBytes, &confE)
	if err != nil {
		return nil, err
	}

	return &confE, nil
}

func loadConfigFromEmbed() (*SensitiveConfig, error) {
	confFile, err := embed.FS.Open("rules.toml")
	if err != nil {
		return nil, err
	}

	confBytes, err := ioutil.ReadAll(confFile)
	if err != nil {
		return nil, err
	}

	confE := SensitiveConfig{}
	err = toml.Unmarshal(confBytes, &confE)
	if err != nil {
		return nil, err
	}

	return &confE, nil
}

func compileGlob() error {
	if sensitiveConfig == nil {
		return errors.New("rules: can't compile glob because sensitiveConfig is nil")
	}

	for _, whitePath := range sensitiveConfig.WhiteList.Paths {
		g, err := glob.Compile(whitePath)
		if err != nil {
			log.Error(err)
		}

		sensitiveConfig.WhiteList.PathsGlob = append(sensitiveConfig.WhiteList.PathsGlob, g)
	}

	return nil
}

func compileRule() error {
	if sensitiveConfig == nil {
		return errors.New("rules: can't compile rule because sensitiveConfig is nil")
	}

	for i, rule := range sensitiveConfig.Rules {
		if rule.Match != "" {
			if strings.HasPrefix(rule.Match, "$contains:") {
				sensitiveConfig.Rules[i].MatchContains = rule.Match
			} else {
				r, err := regexp.Compile(rule.Match)
				if err != nil {
					log.Error(err)
				} else {
					sensitiveConfig.Rules[i].MatchRegex = r
				}
			}
		}

		if rule.Filepath != "" {
			r, err := regexp.Compile(rule.Filepath)
			if err != nil {
				log.Error(err)
				continue
			}

			sensitiveConfig.Rules[i].FilepathRegexp = r
		}
	}

	return nil
}

func SingletonConf() *SensitiveConfig {
	return sensitiveConfig
}

func Init() {
	confFromService, err := loadConfigFromService()
	if err != nil {
		log.Error(err)

		confFromEmbed, err := loadConfigFromEmbed()
		if err != nil {
			log.Error(err)
		} else {
			sensitiveConfig = confFromEmbed
		}
	} else {
		sensitiveConfig = confFromService
	}

	if sensitiveConfig == nil {
		panic(errors.New("rule: can't init sensitiveConfig from service or local"))
	}

	err = compileGlob()
	if err != nil {
		panic(err)
	}

	err = compileRule()
	if err != nil {
		panic(err)
	}
}
