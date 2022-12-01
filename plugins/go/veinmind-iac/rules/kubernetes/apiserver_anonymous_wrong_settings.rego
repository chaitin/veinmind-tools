package brightMirror.kubernetes

import data.common
import future.keywords.in
import future.keywords.contains
import future.keywords.if

risks[res]{
    input.metadata.name=="system:anonymous"
    input.roleRef.name=="cluster-admin"
    res := common.result({"original":input.metadata.name, "Path": input.Path}, "KN-006")
}