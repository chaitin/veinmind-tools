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
    "name": "chwon flag in COPY",
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