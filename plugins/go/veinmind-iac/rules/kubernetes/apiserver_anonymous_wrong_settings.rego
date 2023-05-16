package brightMirror.kubernetes

import data.common

risks[res]{
    input[_].kind=="ClusterRoleBinding"
    input[_].roleRef.name=="cluster-admin"
    input[_].subjects[i].name=="system:anonymous"
    res := common.result({"original":"UnSafeSettings:`metadata.name`,`roleRef.name`", "Path": input[_].Path}, "KN-006")
}