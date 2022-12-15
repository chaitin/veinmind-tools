package common

meta_data["KN-001"] = {
    "id": "KN-001",
    "name": "Pod Container Use Latest Image",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "使用了latest镜像，会导致镜像更新出现非预期的错误",
    "solution": "使用指定的tag镜像版本替代latest",
    "reference": "",
}

meta_data["KN-002"] = {
    "id": "KN-002",
    "name": "Pod Container Process Elevate Self",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "容器中的程序可以提升自己的权限并以root身份运行，这可能会使程序控制容器和节点。",
    "solution": "将 containers[].securityContext.allowPrivilegeEscalation 置为false",
    "reference": "https://kubernetes.io/docs/concepts/security/pod-security-standards/#restricted",
}

meta_data["KN-003"] = {
    "id": "KN-003",
    "name": "Default AppArmor profile not set",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "未设置默认AppArmor配置，容器中的程序可以绕过AppArmor保护策略。",
    "solution": "删除 `container.apparmor.security.beta.kubernetes.io`注释或设置为`runtime/default`。",
    "reference": "https://kubernetes.io/docs/concepts/security/pod-security-standards/#baseline",
}

meta_data["KN-004"] = {
    "id": "KN-004",
    "name": "SElinux set custom options",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "应禁止设置自定义SELinux用户或角色选项。",
    "solution": "删除 `spec.securityContext.seLinuxOptions`、`spec.securityContext.seLinuxOptions` 以及 `spec.initContainers[*].securityContext.seLinuxOptions.`",
    "reference": "https://kubernetes.io/docs/concepts/security/pod-security-standards/#baseline",
}
meta_data["KN-005"] = {
    "id": "KN-005",
    "name": "insecure-port enabled",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "apiserver的insecure-port不会进行身份验证，开启此选项会造成apiserver的未授权访问。",
    "solution": "删除 /etc/kubernetes/manifest/kube-apiserver.yaml中的`--insecure-port=xxx`选项",
    "reference": "",
}
meta_data["KN-006"] = {
    "id": "KN-006",
    "name": "secure-port wrong setting",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "如果错误地将`system:anonymous`用户绑定到`cluster-admin`用户组，则6443 端口允许匿名用户以管理员权限向集群内部下发指令",
    "solution": "删除错误的system:anonymous clusterrolebinding配置",
    "reference": "",
}
meta_data["KN-007"] = {
    "id": "KN-007",
    "name": "10250 port access wrong setting",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "如果在/var/lib/config.yaml中进行了错误的配置，会导致10250端口的未授权访问，通过该端口可以创建pod和控制pod",
    "solution": "将/var/lib/kubelet/config.yml中authentication的anonymous的enabled设置为false，authorization的mode设置为webhook",
    "reference": "",
}
meta_data["KN-008"] = {
    "id": "KN-008",
    "name": "dashboard wrong access",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "如果对dashboard部署采用了错误的配置会导致dashboard登陆界面存在跳过按钮，可以未授权访问",
    "solution": "将dashboard部署的配置文件中的`--enable-skip-login`删除",
    "reference": "",
}
meta_data["KN-009"] = {
    "id": "KN-009",
    "name": "missing --client-cert-auth=true",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "未检测到etcd配置文件中的--client-cert-auth=true选项，可能会造成etcd不安全访问",
    "solution": "在etcd配置文件中添加--client-cert-auth=true配置项",
    "reference": "",
}
meta_data["KN-010"] = {
    "id": "KN-010",
    "name": "missing --peer-client-cert-auth=true",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "未检测到etcd配置文件中的--peer-client-cert-auth=true选项，可能会造成etcd不安全访问",
    "solution": "在etcd配置文件中添加--peer-client-cert-auth=true配置项",
    "reference": "",
}
meta_data["KN-011"] = {
    "id": "KN-011",
    "name": "Pod Container Use Privileged Mode",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "检测到pod内有容器采用特权模式启动，有逃逸风险",
    "solution": "将pod配置文件中的securityContext的privileged设置为false",
    "reference": "",
}
meta_data["KN-012"] = {
    "id": "KN-012",
    "name": "Pod Container Insecure Permission Granted:SYS_ADMIN",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "检测到pod内有容器被授予了SYS_ADMIN权限，拥有该权限的容器可以通过notify_on_release或devices.allow等方式逃逸，有逃逸风险",
    "solution": "将pod配置文件中的capabilities下的SYS_ADMIN权限删除",
    "reference": "",
}
meta_data["KN-013"] = {
    "id": "KN-013",
    "name": "Pod Container Insecure Permission Granted:DAC_READ_SEARCH",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "检测到pod内有容器被授予了CAP_DAC_READ_SEARCH权限，拥有该权限的容器可以通过open_by_handle_at函数读取宿主机的文件，有逃逸风险",
    "solution": "将pod配置文件中的capabilities下的DAC_READ_SEARCH权限删除",
    "reference": "",
}
meta_data["KN-014"] = {
    "id": "KN-014",
    "name": "Pod Container Insecure Permission Granted:SYS_MODULE",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "检测到pod内有容器被授予了SYS_MODULE权限，拥有该权限的容器可以加载内核模块，有逃逸风险",
    "solution": "将pod配置文件中的capabilities下的SYS_MODULE权限删除",
    "reference": "",
}
meta_data["KN-015"] = {
    "id": "KN-015",
    "name": "Pod Container Insecure Permission Granted:DAC_OVERRIDE",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "检测到pod内有容器被授予了DAC_OVERRIDE权限，拥有该权限的容器可以绕过文件读、写、执行权限的检查，有逃逸风险",
    "solution": "将pod配置文件中的capabilities下的DAC_OVERRIDE权限删除",
    "reference": "",
}
meta_data["KN-016"] = {
    "id": "KN-016",
    "name": "Pod Container Insecure File Mount :docker.sock",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "检测到pod内有容器挂载了不安全的目录docker.sock，挂载了docker.sock的容器可以利用向docker.sock发送命令的方式逃逸，有逃逸风险",
    "solution": "将pod配置文件中的volumes下的hostpath挂载的docker.sock删除",
    "reference": "",
}
meta_data["KN-017"] = {
    "id": "KN-017",
    "name": "Pod Container Insecure File Mount :lxcfs",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "检测到pod内有容器挂载了不安全的目录lxcfs，挂载了lxcfs的容器可以利用重写devices.allow等方式进行逃逸，有逃逸风险",
    "solution": "将pod配置文件中的volumes下的hostpath挂载的lxcfs删除",
    "reference": "",
}
meta_data["KN-018"] = {
    "id": "KN-018",
    "name": "Pod Container Insecure File Mount :/",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "检测到pod内有容器挂载了根目录，挂载了根目录的容器可以读写根目录下的文件，有逃逸风险",
    "solution": "将pod配置文件中的volumes下的hostpath挂载的/删除",
    "reference": "",
}
meta_data["KN-019"] = {
    "id": "KN-019",
    "name": "Pod Container Insecure File Mount :/proc",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "检测到pod内有容器挂载了/proc目录，挂载了/proc目录的容器可以通过/proc/sys/kernel/core_pattern的特性在宿主机上执行命令，有逃逸风险",
    "solution": "将pod配置文件中的volumes下的hostpath挂载的/proc删除",
    "reference": "",
}
meta_data["KN-020"] = {
    "id": "KN-020",
    "name": "Pod Container Insecure Permission Granted:SYS_PTRACE",
    "type": "kubernetes",
    "severity": "Medium",
    "description": "检测到pod内有容器被授予了SYS_PTRACE权限，并且pid设置为true。在这种情况下容器可以进行进程代码注入，有逃逸风险",
    "solution": "将pod配置文件中的capabilities下的SYS_PTRACE权限删除,并将hostPid设置为false",
    "reference": "",
}