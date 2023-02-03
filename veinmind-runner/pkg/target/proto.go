package target

import (
	"regexp"
	"strings"
)

type Proto string

const (
	UNKNOWN Proto = "unknown"
	// DOCKERD CONTAINERD: container protol
	// DOCKERD CONTAINERD REGISTRY TARBALL: image protol
	DOCKERD        Proto = "dockerd"
	CONTAINERD     Proto = "containerd"
	REGISTRY       Proto = "registry"
	REGISTRY_IMAGE Proto = "registry-image"
	TARBALL        Proto = "tarball"

	// LOCAL GIT KUBERNETES: iac protol
	LOCAL      Proto = "host"
	GIT        Proto = "git"
	KUBERNETES Proto = "kubernetes"
)

func (p Proto) String() string {
	return string(p)
}

func IsProto(t string) bool {
	switch t {
	case DOCKERD.String(), CONTAINERD.String(), REGISTRY.String(), REGISTRY_IMAGE.String(), LOCAL.String(), GIT.String(), KUBERNETES.String(), TARBALL.String():
		return true
	default:
		return false
	}
}

func ParseProto(cmd string, arg string) (Proto, string) {
	switch cmd {
	case "image":
		return regexParse(arg, DOCKERD, generatePattern(DOCKERD, CONTAINERD, REGISTRY, REGISTRY_IMAGE, TARBALL))
	case "container":
		return regexParse(arg, DOCKERD, generatePattern(DOCKERD, CONTAINERD))
	case "iac":
		return regexParse(arg, LOCAL, generatePattern(LOCAL, GIT, KUBERNETES))
	default:
		return UNKNOWN, ""
	}
}

func regexParse(arg string, defaultProto Proto, patternString string) (Proto, string) {
	pattern := regexp.MustCompile(patternString)
	matches := pattern.FindStringSubmatch(arg)

	// complete matches
	if len(matches) == 3 {
		if IsProto(strings.Trim(matches[1], ":")) {
			return Proto(strings.Trim(matches[1], ":")), matches[2]
		} else {
			// means user input is value
			// use default proto
			return defaultProto, matches[2]
		}
	}

	// no match
	// use default proto and empty value
	return defaultProto, ""
}

func generatePattern(ps ...Proto) string {
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
