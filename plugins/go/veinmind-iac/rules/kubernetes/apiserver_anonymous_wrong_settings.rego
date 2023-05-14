package brightMirror.kubernetes

import data.common

risks[res]{
    input.kind=="ClusterRoleBinding"
    input.roleRef.name=="cluster-admin"
    input.subjects[i].name=="system:anonymous"
    res := common.result({"original":"UnSafeSettings:`metadata.name`,`roleRef.name`", "Path": input.Path}, "KN-006")
}