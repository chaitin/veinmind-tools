package common

meta_data["KN-001"] = {
    "id": "KN-001",
    "name": "Pod Container Use Latest Image",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "Using the latest image will cause an unexpected error in the image update",
    "solution": "Use the specified tag image version instead of latest",
    "reference": "",
}

meta_data["KN-002"] = {
    "id": "KN-002",
    "name": "Pod Container Process Elevate Self",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "Programs in a container can elevate their privileges and run as root, potentially giving them control over containers and nodes.",
    "solution": "Set containers[].SecurityContext.AllowPrivilegeEscalation to false",
    "reference": "https://kubernetes.io/docs/concepts/security/pod-security-standards/#restricted",
}

meta_data["KN-003"] = {
    "id": "KN-003",
    "name": "Default AppArmor profile not set",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "Without the default AppArmor configuration set, programs in the container can bypass the AppArmor protection policy.",
    "solution": "Delete `container.Apparmor.Security.Beta.Kubernetes.IO` comments or set to `runtime/default`.",
    "reference": "https://kubernetes.io/docs/concepts/security/pod-security-standards/#baseline",
}

meta_data["KN-004"] = {
    "id": "KN-004",
    "name": "SELinux set custom options",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "Setting custom SELinux user or role options should be disallowed.",
    "solution": "Delete `spec.securityContext.seLinuxOptions`、`spec.securityContext.seLinuxOptions` and `spec.initContainers[*].securityContext.seLinuxOptions.`",
    "reference": "https://kubernetes.io/docs/concepts/security/pod-security-standards/#baseline",
}
meta_data["KN-005"] = {
    "id": "KN-005",
    "name": "insecure-port enabled",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "Insecure-port of apiserver does not perform authentication. Enabling this option will cause unauthorized access to apiserver.",
    "solution": "Delete the option `--insecure-port=xxx` in /etc/kubernetes/manifest/kube-apiserver.yaml",
    "reference": "",
}
meta_data["KN-006"] = {
    "id": "KN-006",
    "name": "secure-port wrong setting",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "If a 'system:anonymous' user is mistakenly bound to a' cluster-admin 'user group, port 6443 allows anonymous users to issue commands internally to the cluster with administrator rights",
    "solution": "Delete the setting of system:anonymous clusterrolebinding",
    "reference": "",
}
meta_data["KN-007"] = {
    "id": "KN-007",
    "name": "10250 port access wrong setting",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "Incorrect configuration in /var/lib/config.yaml can result in unauthorized access to port 10250, through which Pods can be created and controlled",
    "solution": "Set `authentication:anonymous:enabled` to false and `authorization:mode` to `webhook` in /var/lib/kubelet/config.yml",
    "reference": "",
}
meta_data["KN-008"] = {
    "id": "KN-008",
    "name": "dashboard wrong access",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "If you configure the dashboard deployment incorrectly, the skip button may exist on the dashboard login page, and unauthorized access is possible",
    "solution": "Delete `--enable-skip-login` in dashboard configure file",
    "reference": "",
}
meta_data["KN-009"] = {
    "id": "KN-009",
    "name": "missing --client-cert-auth=true",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "The `--client-cert-auth=true` option in the etcd configuration file is not detected, which may result in insecure etcd access",
    "solution": "Add the ``--client-cert-auth=true` configuration item to the etcd configuration file",
    "reference": "",
}
meta_data["KN-010"] = {
    "id": "KN-010",
    "name": "missing --peer-client-cert-auth=true",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "The `--peer-client-cert-auth=true` option in the etcd configuration file is not detected, which may result in insecure etcd access",
    "solution": "Add the ``--peer-client-cert-auth=true` configuration item to the etcd configuration file",
    "reference": "",
}
meta_data["KN-011"] = {
    "id": "KN-011",
    "name": "Pod Container Use Privileged Mode",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "A container in pod has been started in privileged mode, and there is an escape risk. Procedure",
    "solution": "Set the privileged of the securityContext in the pod profile to false",
    "reference": "",
}
meta_data["KN-012"] = {
    "id": "KN-012",
    "name": "Pod Container Insecure Permission Granted:SYS_ADMIN",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "It is detected that a container in pod has been granted SYS_ADMIN permission. Containers with this permission can escape by notify_on_release or devices.allow",
    "solution": "Delete the SYS_ADMIN permissions in capabilities in the pod configuration file",
    "reference": "",
}
meta_data["KN-013"] = {
    "id": "KN-013",
    "name": "Pod Container Insecure Permission Granted:DAC_READ_SEARCH",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "It is detected that a container in pod has the CAP_DAC_READ_SEARCH permission. The container with this permission can read files of the host through the open_by_handle_at function, which is at risk of escaping. Procedure",
    "solution": "Delete the CAP_DAC_READ_SEARCH permissions in capabilities in the pod configuration file",
    "reference": "",
}
meta_data["KN-014"] = {
    "id": "KN-014",
    "name": "Pod Container Insecure Permission Granted:SYS_MODULE",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "A container in pod has been granted SYS_MODULE permission. A container with this permission can load kernel modules and is at risk of escaping",
    "solution": "Delete the SYS_MODULE permission in capabilities in the pod configuration file",
    "reference": "",
}
meta_data["KN-015"] = {
    "id": "KN-015",
    "name": "Pod Container Insecure Permission Granted:DAC_OVERRIDE",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "It is detected that a container in pod is granted DAC_OVERRIDE permission. A container with this permission can bypass the check of file read, write, and execute permissions, causing risks to escape",
    "solution": "Delete the DAC_OVERRIDE permission in capabilities in the pod configuration file",
    "reference": "",
}
meta_data["KN-016"] = {
    "id": "KN-016",
    "name": "Pod Container Insecure File Mount :docker.sock",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "Unsafe directory docker.sock is mounted to a container in pod. Containers mounted with docker.sock can escape by sending commands to docker.sock, and there is a risk of escape",
    "solution": "Delete `docker.sock` attached to hostpath on volumes in the pod configuration file",
    "reference": "",
}
meta_data["KN-017"] = {
    "id": "KN-017",
    "name": "Pod Container Insecure File Mount :lxcfs",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "An insecure directory lxcfs is mounted to a container in pod. The container mounted with lxcfs can escape by means of overriding devices.allow",
    "solution": "Delete `lxcfs` that are attached to hostpath of volumes in the pod configuration file",
    "reference": "",
}
meta_data["KN-018"] = {
    "id": "KN-018",
    "name": "Pod Container Insecure File Mount :/",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "The root directory is mounted to a container in pod. The container mounted to the root directory can read and write files in the root directory, causing an escape risk. Procedure",
    "solution": "Delete `/` that are attached to hostpath of volumes in the pod configuration file",
    "reference": "",
}
meta_data["KN-019"] = {
    "id": "KN-019",
    "name": "Pod Container Insecure File Mount :/proc",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "The /proc directory is mounted to a container in pod. The container can run commands on the host using the /proc/sys/kernel/core_pattern feature, which has a risk of escaping.",
    "solution": "Delete `/proc` that are attached to hostpath of volumes in the pod configuration file",
    "reference": "",
}
meta_data["KN-020"] = {
    "id": "KN-020",
    "name": "Pod Container Insecure Permission Granted:SYS_PTRACE",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "Detected that a container in pod has been granted SYS_PTRACE permission and pid is set to true. In this case, the container can do process code injection, with the risk of escape",
    "solution": "Delete the DAC_OVERRIDE permission in capabilities in the pod configuration file and set hostPid to false",
    "reference": "",
}