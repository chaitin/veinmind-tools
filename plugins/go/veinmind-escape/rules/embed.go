package rules

import (
	"embed"
	"io/fs"
	"os"
)

//go:embed rule.toml
var RuleFS embed.FS

func Open(name string) (fs.File, error) {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return RuleFS.Open(name)
	} else {
		return os.Open(name)
	}
}
