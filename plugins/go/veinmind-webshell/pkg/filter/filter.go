// Package filter find web script from image/container
package filter

import (
	"io/fs"
)

var Kit *kit

type kit struct{}

// Filter indicates whether the file is a web script
func (f kit) Filter(path string, info fs.FileInfo) (bool, ScriptType, error) {
	for _, suffix := range scriptSuffixes {
		if suffix.Match(info.Name()) {
			if t, ok := scriptSuffixTypeMap[suffix]; ok {
				return true, t, nil
			}
		}
	}

	return false, UNKNOWN_TYPE, nil
}

func init() {
	Kit = new(kit)
}
