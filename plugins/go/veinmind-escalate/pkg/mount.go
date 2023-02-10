package pkg

import (
	"github.com/chaitin/veinmind-common-go/service/report/event"
	"path/filepath"

	api "github.com/chaitin/libveinmind/go"
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

func ContainerUnsafeMount(fs api.FileSystem) ([]*event.EscapeDetail, error) {
	var res = make([]*event.EscapeDetail, 0)
	container := fs.(api.Container)
	spec, err := container.OCISpec()
	if err != nil {
		return res, err
	}

	for _, mount := range spec.Mounts {
		for _, pattern := range UnsafeMountPaths {
			matched, _ := filepath.Match(pattern, mount.Source)
			if matched {
				// /var/log逃逸仅在k8s环境下
				if pattern == "/var/log" {
					// 不存在/var/run/secrets/kubernetes.io/serviceaccount/token, 则非k8s容器，跳过。
					if _, err := fs.Stat("/var/run/secrets/kubernetes.io/serviceaccount/token"); err != nil {
						continue
					}
				}
				res = append(res, &event.EscapeDetail{
					Target: mount.Source,
					Reason: MOUNTREASON,
					Detail: "UnSafeMountDestination " + mount.Destination,
				})
			}
		}
	}
	return res, nil
}

func init() {
	ContainerCheckList = append(ContainerCheckList, ContainerUnsafeMount)
}
