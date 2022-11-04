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