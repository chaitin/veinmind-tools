package rules

import (
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/pelletier/go-toml"
)

type Exp struct {
	Exp string `toml:"exp"`
}
type Exps struct {
	Shell                      []*Exp `toml:"shell"`
	Command                    []*Exp `toml:"command"`
	ReverseShell               []*Exp `toml:"reverse-shell"`
	NonInteractiveReverseShell []*Exp `toml:"non-interactive-reverse-shell"`
	BindShell                  []*Exp `toml:"bind-shell"`
	NonInteractiveBindShell    []*Exp `toml:"non-interactive-bind-shell"`
	FileUpload                 []*Exp `toml:"file-upload"`
	FileDownload               []*Exp `toml:"file-download"`
	FileWrite                  []*Exp `toml:"file-write"`
	FileRead                   []*Exp `toml:"file-read"`
	LibraryLoad                []*Exp `toml:"library-load"`
	SUID                       []*Exp `toml:"suid"`
	Sudo                       []*Exp `toml:"sudo"`
	Capabilities               []*Exp `toml:"capabilities"`
	LimitedSUID                []*Exp `toml:"limited-suid"`
}

type Rule struct {
	Name        string   `toml:"Name"`
	Description string   `toml:"Description"`
	Tags        []string `toml:"Tags"`
	Exps        Exps     `toml:"exps"`
}

type Config struct {
	Rules []*Rule `toml:"privilege-esclation"`
}

func GetRuleFromFile() (*Config, error) {
	var config Config
	content, err := Readfile("rule.toml")
	if err != nil {
		log.Fatal(err)
	}

	if err := toml.Unmarshal(content, &config); err != nil {
		log.Error(err)
		return nil, err
	}

	return &config, nil
}
