package brightMirror.kubernetes

import data.common
import future.keywords.every
import future.keywords.in
import future.keywords.contains
import future.keywords.if

risks[res]{
    input.spec.containers[i].securityContext.privileged==true
        d := input.spec.containers[i].name
        a=["UnsafeContainers"]
        b=[d]
        c:=array.concat(a,b)
        res := common.result({"original":concat(":",c), "Path": input.Path}, "KN-011")

}