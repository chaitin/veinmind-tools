package target

import (
	"regexp"
	"strings"
)

type Protol string

const (
	UNKNOWN Protol = "unknown"
	// DOCKERD CONTAINERD: container protol
	// DOCKERD CONTAINERD REGISTRY TARBALL: image protol
	DOCKERD    Protol = "dockerd"
	CONTAINERD Protol = "containerd"
	REGISTRY   Protol = "registry"
	TARBALL    Protol = "tarball"

	// LOCAL GIT KUBERNETES: iac protol
	LOCAL      Protol = "host"
	GIT        Protol = "git"
	KUBERNETES Protol = "kubernetes"
)

func (p Protol) String() string {
	return string(p)
}

func IsProto(t string) bool {
	switch t {
	case DOCKERD.String(), CONTAINERD.String(), REGISTRY.String(), LOCAL.String(), GIT.String(), KUBERNETES.String():
		return true
	default:
		return false
	}
}

func ParseProto(cmd string, arg string) (Protol, string) {
	switch cmd {
	case "image":
		return regexParse(arg, DOCKERD, generatePattern(DOCKERD, CONTAINERD, REGISTRY))
	case "container":
		return regexParse(arg, DOCKERD, generatePattern(DOCKERD, CONTAINERD))
	case "iac":
		return regexParse(arg, LOCAL, generatePattern(LOCAL, GIT, KUBERNETES))
	default:
		return UNKNOWN, ""
	}
}

func regexParse(arg string, defaultProtol Protol, patternString string) (Protol, string) {
	pattern := regexp.MustCompile(patternString)
	matches := pattern.FindStringSubmatch(arg)

	// complete matches
	if len(matches) == 3 {
		if IsProto(strings.Trim(matches[1], ":")) {
			return Protol(strings.Trim(matches[1], ":")), matches[2]
		} else {
			// means user input is value
			// use default protol
			return defaultProtol, matches[2]
		}
	}
	// bad matches
	return UNKNOWN, ""
}

func generatePattern(ps ...Protol) string {
	pattern := ""
	for i, p := range ps {
		if i == 0 {
			pattern = "("
		}
		pattern += p.String() + ":"
		if i == len(ps)-1 {
			pattern += ")"
		} else {
			pattern += "|"
		}
	}
	if pattern != "" {
		pattern += "?(.*)"
	}

	return pattern
}
