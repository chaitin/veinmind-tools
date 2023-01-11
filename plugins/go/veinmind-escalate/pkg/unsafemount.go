package pkg

import (
	api "github.com/chaitin/libveinmind/go"
	"path/filepath"
)

var UnsafeMountPaths = []string{
	"/lxcfs",
	"/",
	"/root",
	"/etc",
	"/var",
	"/proc",
	"/bin",
	"/sys",
	"/var/run/docker.sock",
	"/run/containerd.sock",
	"/var/run/crio/crio.sock",
	"/var/lib/kubelet",
	"/var/lib/kubelet/pki",
	"/etc/kubernetes",
	"/etc/kubernetes/manifests",
}

func DetectContainerUnsafeMount(fs api.FileSystem) error {
	container := fs.(api.Container)
	spec, err := container.OCISpec()
	if err != nil {
		return err
	}

	for _, mount := range spec.Mounts {
		for _, pattern := range UnsafeMountPaths {
			matched, err := filepath.Match(pattern, mount.Source)
			if err != nil {
				continue
			}

			if matched {
				AddResult(mount.Destination, MOUNTREASON, "UnSafeMount "+mount.Destination)
			}
		}
	}

	return nil
}
