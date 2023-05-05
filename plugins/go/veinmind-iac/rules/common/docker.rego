package common

meta_data["DF-001"] = {
   "id": "DF-001",
   "name": "expose port 22",
   "type": "dockerfile",
   "severity": "Medium",
   "description": "The container with port 22 open may allow users to log in through SSH. Therefore, check whether the application needs to use port 22 and avoid using port 22",
   "solution": "Cancel the Expose operation on port 22",
   "reference": "",
}

meta_data["DF-002"]  := {
    "id": "DF-002",
    "name": "expose out of range [1-65535]",
    "type": "dockerfile",
    "severity": "High",
    "description": "The Expose port number is out of the normal range of ports: [0-65535]",
    "solution": "Make sure the Expose port number is correct",
    "reference": "",
}

meta_data["DF-003"] := {
    "id": "DF-003",
    "name": "from latest image",
    "type": "dockerfile",
    "severity": "Medium",
    "description": "Using the latest image will cause an unexpected error in the image update",
    "solution": "Use the specified tag image version instead of latest",
    "reference": "",
}

meta_data["DF-004"] := {
    "id": "DF-004",
    "name": "from use platform",
    "type": "dockerfile",
    "severity": "Low",
    "description": "The platform parameter was found in the FROM directive. Please try to remove unnecessary platform designations",
    "solution": "Removes the platform designation",
    "reference": "",
}

meta_data["DF-005"] := {
    "id": "DF-005",
    "name": "user root",
    "type": "dockerfile",
    "severity": "High",
    "description": "Starting the container with root can cause container escape. The best practice is to start the container with 'USER' in the Dockerfile as a non-root user",
    "solution": "Add 'USER xxxx' to the Dockerfile and make sure the user is not root",
    "reference": "",
}

meta_data["DF-006"] := {
    "id": "DF-006",
    "name": "WORKDIR not absolute",
    "type": "dockerfile",
    "severity": "High",
    "description": "WORKDIR uses a non-absolute path",
    "solution": "Use absolute paths in WORKDIR",
    "reference": "https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#workdir",
}


meta_data["DF-007"] := {
    "id": "DF-007",
    "name": "using HEALTHCHECK",
    "type": "dockerfile",
    "severity": "Low",
    "description": "Use HEALTHCHECK to ensure normal services",
    "solution": "Use HEALTHCHECK to keep your business running",
    "reference": "https://docs.docker.com/engine/reference/builder/#healthcheck",
}

meta_data["DF-008"] := {
    "id": "DF-008",
    "name": "multiple HEALTHCHECK",
    "type": "dockerfile",
    "severity": "Medium",
    "description": "HEALTHCHECK is repeated",
    "solution": "Keep the unique HEALTHCHECK",
    "reference": "https://docs.docker.com/engine/reference/builder/#healthcheck",
}

meta_data["DF-009"] := {
    "id": "DF-009",
    "name": "chown flag in COPY",
    "type": "dockerfile",
    "severity": "Low",
    "description": "Make sure not to use the --chown parameter when the user only needs enforcement rights",
    "solution": "delete --chown",
    "reference": "",
}

meta_data["DF-010"] := {
    "id": "DF-010",
    "name": "use COPY instead of ADD",
    "type": "dockerfile",
    "severity": "Low",
    "description": "Use COPY instead of ADD unless you want to extract a tar file. The ADD directive extracts tar files, which increases the risk of ZIP-based vulnerabilities. Therefore, the COPY directive, which does not extract tar files, is recommended",
    "solution": "Use COPY instead of the ADD directive",
    "reference": "https://docs.docker.com/engine/reference/builder/#add",
}

meta_data["DF-011"] := {
    "id": "DF-011",
    "name": "use ADD to fetch package",
    "type": "dockerfile",
    "severity": "Low",
    "description": "Retrieving packages remotely using ADD is extremely dangerous. Use curl/wget to retrieve packages",
    "solution": "Use curl/wget instead of ADD",
    "reference": "https://docs.docker.com/engine/reference/builder/#add",
}

meta_data["DF-012"] := {
    "id": "DF-012",
    "name": "secret value in ENV",
    "type": "dockerfile",
    "severity": "Medium",
    "description": "Sensitive information is written to ENV, which may leak",
    "solution": "Remove unless necessary sensitive information",
    "reference": "",
}

meta_data["DF-013"] := {
    "id": "DF-013",
    "name": "both use wget and curl",
    "type": "dockerfile",
    "severity": "Low",
    "description": "Both wget and curl are used, and they do the same thing",
    "solution": "Make sure you use the same tool that does the same thing",
    "reference": "https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#run",
}

meta_data["DF-014"] := {
    "id": "DF-014",
    "name": "run with sudo",
    "type": "dockerfile",
    "severity": "CRITICAL",
    "description": "Avoid using 'sudo' in the 'RUN' instruction, which may lead to unexpected results",
    "solution": "Avoid using sudo",
    "reference": "https://docs.docker.com/engine/reference/builder/#run",
}

meta_data["DF-015"] := {
    "id": "DF-015",
    "name": "add with parent directory",
    "type": "dockerfile",
    "severity": "Low",
    "description": "Avoid using ADD with parent directory",
    "solution": "Change path to absolute path",
    "reference": "https://docs.docker.com/engine/reference/builder/#add",
}

meta_data["DF-016"] := {
    "id": "DF-016",
    "name": "Multiple CMD commands are used",
    "type": "dockerfile",
    "severity": "CRITICAL",
    "description": "There can only be one CMD command in a Docker file. If you list multiple CMDs, only the last one will take effect.",
    "solution": "Delete other cmd commands and use only one cmd command",
    "reference": "https://docs.docker.com/engine/reference/builder/#cmd",
}

meta_data["DF-017"] := {
    "id": "DF-017",
    "name": "Multiple entrypoint commands are used",
    "type": "dockerfile",
    "severity": "CRITICAL",
    "description": "There can only be one entrypoint command in a Docker file. If you list multiple entrypoints, only the last one will take effect.",
    "solution": "Delete other entrypoint commands and use only one entrypoint command",
    "reference": "https://docs.docker.com/engine/reference/builder/#entrypoint",
}
meta_data["DF-018"] := {
    "id": "DF-018",
    "name": "Duplicate aliases defined in different FROMs",
    "type": "dockerfile",
    "severity": "CRITICAL",
    "description": "Different FROMs cannot have the same alias definition.",
    "solution": "Change the alias to be different",
    "reference": "https://docs.docker.com/develop/develop-images/multistage-build/",
}
meta_data["DF-019"] := {
    "id": "DF-019",
    "name": "`COPY --from` comes from the current image",
    "type": "dockerfile",
    "severity": "MEDIUM",
    "description": "COPY --from` should not use the from alias of the current mirror, as it cannot copy from itself.",
    "solution": "Change `--from` so it doesn't refer to itself",
    "reference": "https://docs.docker.com/develop/develop-images/multistage-build/",
}
meta_data["DF-020"] := {
    "id": "DF-020",
    "name": "`COPY --from` does not use the `--link` argument",
    "type": "dockerfile",
    "severity": "MEDIUM",
    "description": "Use the `--link` parameter to avoid the need to rebuild the intermediate stages when the build fails. `--link` will reuse the previously generated layer and merge it on top of the new layer. When the base image receives an update, you can Easily rebase an image without having to do the entire build again.",
    "solution": "add --link arg",
    "reference": " https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#leverage-build-cache",
}
meta_data["DF-021"] := {
    "id": "DF-021",
    "name": "Using `FROM` specific `ARG` variable",
    "type": "dockerfile",
    "severity": "MEDIUM",
    "description": "The `ARG` variable before `FROM` only applies to `FROM`, and cannot be obtained in subsequent instructions. To use it, you need to redefine `ARG`",
    "solution": "Please redefine `ARG` after `FROM`",
    "reference": "https://docs.docker.com/engine/reference/builder/#understand-how-arg-and-from-interact",
}
meta_data["DF-022"] := {
    "id": "DF-022",
    "name": "Use of deprecated directive MAINTAINER",
    "type": "dockerfile",
    "severity": "High",
    "description": "MAINTAINER is deprecated since Docker 1.13.0.",
    "solution": "Use LABEL instead of MAINTAINER",
    "reference": "https://docs.docker.com/engine/deprecated/#maintainer-in-dockerfile",
}

meta_data["DF-023"] := {
    "id": "DF-023",
    "name": "Complicated `RUN cd` syntax is used",
    "type": "dockerfile",
    "severity": "MEDIUM",
    "description": "You should use `WORKDIR` instead of a plethora of instructions like `RUN cd…`, which are difficult to read, troubleshoot and maintain",
    "solution": "use `WORKDIR` instead of complicated `RUN cd` switching working directory",
    "reference": "https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#workdir",
}

meta_data["DF-024"] := {
    "id": "DF-024",
    "name": "Used `RUN <package-manager> update` alone",
    "type": "dockerfile",
    "severity": "MEDIUM",
    "description": "The directive `RUN <package-manager> update` should always be followed by `<package-manager> install` within the same RUN statement.",
    "solution": "Merge `<packagemanager>update` and `<package manager> install` directives into one directive",
    "reference": "https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#run",
}

meta_data["DF-025"] := {
    "id": "DF-025",
    "name": "RUN <package-manager> install` does not use the `-y` parameter",
    "type": "dockerfile",
    "severity": "MEDIUM",
    "description": "The `-y` parameter should be used after `RUN <package-manager> install` to avoid manual input.",
    "solution": "Add `-y` such as: `yum install -y`, `apt-get install -y`, ``",
    "reference": "https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#run",
}

meta_data["DF-026"] := {
    "id": "DF-026",
    "name": "Didn't clean up after using `RUN apk add`",
    "type": "dockerfile",
    "severity": "MEDIUM",
    "description": "Cleanup command should be used after `RUN apk add` to clear package cache data and reduce image size.",
    "solution": "Add a cleanup command such as: `apk add --no-cache`",
    "reference": "https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#run",
}

meta_data["DF-027"] := {
    "id": "DF-027",
    "name": "Didn't clean up after using `RUN dnf install`",
    "type": "dockerfile",
    "severity": "MEDIUM",
    "description": "Cleanup command should be used after `RUN dnf install` to clear package cache data and reduce image size.",
    "solution": "Add clean command like: `dnf clean all`",
    "reference": "https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#run",
}

meta_data["DF-028"] := {
    "id": "DF-028",
    "name": "Didn't clean up after using `RUN yum install`",
    "type": "dockerfile",
    "severity": "MEDIUM",
    "description": "Cleanup command should be used after `RUN yum install` to clear package cache data and reduce image size.",
    "solution": "Add cleaning commands such as: `yum clean all`",
    "reference": "https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#run",
}

meta_data["DF-029"] := {
    "id": "DF-029",
    "name": "Didn't clean up after using `RUN zypper in/remove/source-install`",
    "type": "dockerfile",
    "severity": "MEDIUM",
    "description": "Cleanup command should be used after `RUN apk add` to clear package cache data and reduce image size.",
    "solution": "Add cleaning commands such as: `zypper clean` or `zypper cc`",
    "reference": "https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#run",
}

meta_data["DF-030"] := {
    "id": "DF-030",
    "name": "RUN useradd` needs to add `--no-log-init` parameter",
    "type": "dockerfile",
    "severity": "MEDIUM",
    "description": "Due to an unresolved bug in the Go archive/tar package's handling of sparse files, attempting to create a user with a very large UID in a Docker container could lead to disk exhaustion as /var/log/faillog in the container layer fills up Empty characters. The workaround is to pass the --no-log-init flag to useradd. The Debian/Ubuntu adduser wrapper does not support this flag.",
    "solution": "Use the `--no-log-init` parameter in `useradd`",
    "reference": "https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#run",
}