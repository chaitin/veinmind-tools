package pkg

import (
	api "github.com/chaitin/libveinmind/go"
	"path/filepath"
)

var UnsafeMountPaths = []string{
	"/lxcfs",
	"/",
	"/etc",
	"/var",
	"/proc",
	"/sys",
	"/etc/crontab",
	"/etc/passwd",
	"/etc/shadow",
	"/root/.ssh",

	"/var/run/docker.sock",
	"/run/containerd.sock",
	"/var/run/crio/crio.sock",

	"/var/lib/kubelet",
	"/var/lib/kubelet/pki",
	"/etc/kubernetes",
	"/etc/kubernetes/manifests",
	"/var/log",
}

func ContainerUnsafeMount(fs api.FileSystem) error {
	container := fs.(api.Container)
	spec, err := container.OCISpec()
	if err != nil {
		return err
	}

	for _, mount := range spec.Mounts {
		for _, pattern := range UnsafeMountPaths {
			var matched bool
			if pattern == "/var/log" {
				content, errTOKEN := fs.Open("/var/run/secrets/kubernetes.io/serviceaccount/token") //检查是否在k8s pod环境下
				if errTOKEN == nil {
					matched, err = filepath.Match(pattern, mount.Source)
					if err != nil {
						continue
					}
				}
				defer FileClose(content, errTOKEN)
			} else {
				matched, err = filepath.Match(pattern, mount.Source)
				if err != nil {
					continue
				}
			}

			if matched {
				AddResult(mount.Source, MOUNTREASON, "UnSafeMountDestination "+mount.Destination)
			}
		}
	}
	return nil
}
