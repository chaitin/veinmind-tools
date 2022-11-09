package common

meta_data["DF-001"] = {
   "id": "DF-001",
   "name": "expose port 22",
   "type": "dockerfile",
   "severity": "Medium",
   "description": "开放22端口的容器可能允许用户通过SSH进行登录，请确认应用是否需要使用，并避免使用22端口",
   "solution": "取消22端口的Expose操作",
   "reference": "",
}

meta_data["DF-002"]  := {
    "id": "DF-002",
    "name": "expose out of range [1-65535]",
    "type": "dockerfile",
    "severity": "High",
    "description": "Expose 端口号超出了端口的正常范围: [0-65535]",
    "solution": "确认 Expose 端口号正确",
    "reference": "",
}

meta_data["DF-003"] := {
    "id": "DF-003",
    "name": "from latest image",
    "type": "dockerfile",
    "severity": "Medium",
    "description": "使用了latest镜像，会导致镜像更新出现非预期的错误",
    "solution": "使用指定的tag镜像版本替代latest",
    "reference": "",
}

meta_data["DF-004"] := {
    "id": "DF-004",
    "name": "from use platform",
    "type": "dockerfile",
    "severity": "Low",
    "description": "在FROM指令发现了platform参数，请尽量去除不必要的platform指定",
    "solution": "去除platform的指定",
    "reference": "",
}

meta_data["DF-005"] := {
    "id": "DF-005",
    "name": "user root",
    "type": "dockerfile",
    "severity": "High",
    "description": "使用root启动容器可能会造成容器逃逸，最佳实践是在Dockerfile通过`USER`指定非root用户启动",
    "solution": "添加 `USER xxxx` 在Dockerfile中并确保用户为非root用户",
    "reference": "",
}

meta_data["DF-006"] := {
    "id": "DF-006",
    "name": "WORKDIR not absolute",
    "type": "dockerfile",
    "severity": "High",
    "description": "WORKDIR 使用了非绝对路径",
    "solution": "在WORKDIR中使用绝对路径",
    "reference": "https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#workdir",
}


meta_data["DF-007"] := {
    "id": "DF-007",
    "name": "using HEALTHCHECK",
    "type": "dockerfile",
    "severity": "Low",
    "description": "请使用HEALTHCHECK保证业务正常",
    "solution": "使用HEALTHCHECK来保证业务正常运行",
    "reference": "https://docs.docker.com/engine/reference/builder/#healthcheck",
}

meta_data["DF-008"] := {
    "id": "DF-008",
    "name": "multiple HEALTHCHECK",
    "type": "dockerfile",
    "severity": "Medium",
    "description": "重复声明了 HEALTHCHECK",
    "solution": "保留唯一的 HEALTHCHECK",
    "reference": "https://docs.docker.com/engine/reference/builder/#healthcheck",
}

meta_data["DF-009"] := {
    "id": "DF-009",
    "name": "chwon flag in COPY",
    "type": "dockerfile",
    "severity": "Low",
    "description": "当用户只需要执行权限时，确保不使用--chown参数",
    "solution": "去除 --chown",
    "reference": "",
}


meta_data["DF-010"] := {
    "id": "DF-010",
    "name": "use COPY instead of ADD",
    "type": "dockerfile",
    "severity": "Low",
    "description": "除非要提取tar文件，否则应使用COPY而不是ADD。 ADD指令将提取tar文件，这增加了基于Zip的漏洞的风险。因此，建议使用不提取tar文件的COPY指令",
    "solution": "使用COPY代替ADD指令",
    "reference": "https://docs.docker.com/engine/reference/builder/#add",
}

meta_data["DF-011"] := {
    "id": "DF-011",
    "name": "use ADD to fetch package",
    "type": "dockerfile",
    "severity": "Low",
    "description": "使用ADD从远程获取packages是极度危险的，请使用curl/wget来获取packages",
    "solution": "使用curl/wget来替代ADD",
    "reference": "https://docs.docker.com/engine/reference/builder/#add",
}

meta_data["DF-012"] := {
    "id": "DF-012",
    "name": "secret value in ENV",
    "type": "dockerfile",
    "severity": "Medium",
    "description": "ENV内写入了敏感信息，可能会导致信息泄漏",
    "solution": "去除非必要的敏感信息",
    "reference": "",
}

meta_data["DF-013"] := {
    "id": "DF-013",
    "name": "both use wget and curl",
    "type": "dockerfile",
    "severity": "Low",
    "description": "同时使用了wget和curl，两者的功能相同",
    "solution": "确保使用同一种相同功能的工具",
    "reference": "https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#run",
}

meta_data["DF-014"] := {
    "id": "DF-014",
    "name": "run with sudo",
    "type": "dockerfile",
    "severity": "CRITICAL",
    "description": "避免使用`sudo` 在`RUN`指令中，这可能会导致非预期的结果",
    "solution": "避免使用sudo",
    "reference": "https://docs.docker.com/engine/reference/builder/#run",
}