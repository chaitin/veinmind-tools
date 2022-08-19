package filter

import (
	"strings"
)

func (s ScriptSuffix) String() string {
	return string(s)
}

func (s ScriptSuffix) Match(name string) bool {
	return strings.HasSuffix(name, s.String())
}

func (s ScriptType) String() string {
	return string(s)
}
