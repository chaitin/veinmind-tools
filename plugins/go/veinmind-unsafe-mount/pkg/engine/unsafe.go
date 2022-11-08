package engine

var UnsafeMountPaths = []string{
	// system directory
	"/",
	"/root",
	"/etc",
	"/boot",
	"/var",
	"/proc",
	"/bin",
	"/sys",

	// runtime socket
	"/var/run/docker.sock",
	"/run/containerd.sock",
	"/var/run/crio/crio.sock",

	// kubernetes
	"/var/lib/kubelet",
	"/var/lib/kubelet/pki",
	"/etc/kubernetes",
	"/etc/kubernetes/manifests",
}
