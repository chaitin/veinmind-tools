package capability

import (
	api "github.com/chaitin/libveinmind/go"
	"github.com/syndtr/gocapability/capability"
)

func IsPrivileged(container api.Container) bool {
	ocispec, err := container.OCISpec()
	if err != nil {
		return false
	}
	if ocispec.Process == nil {
		return false
	}

	// construct cap map
	capsEffective := ocispec.Process.Capabilities.Effective
	capsEffectiveMap := make(map[string]bool)
	for _, c := range capsEffective {
		capsEffectiveMap[c] = true
	}

	// round all capability
	allCaps := capability.List()
	isPrivileged := true
	for _, c := range allCaps {
		if _, ok := capsEffectiveMap[c.String()]; ok {
			continue
		} else {
			isPrivileged = false
			break
		}
	}

	return isPrivileged
}
