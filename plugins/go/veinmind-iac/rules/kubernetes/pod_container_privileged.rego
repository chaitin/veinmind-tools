package brightMirror.kubernetes

import data.common
import future.keywords.every
import future.keywords.in
import future.keywords.contains
import future.keywords.if

unsafe_privileged[d]{
	input.spec.containers[i].securityContext.privileged==true
    d := input.spec.containers[i].name
}
unSafe_privileged:={
    "UnSafeContainersName":unsafe_privileged
}
risks[res]{
    count(unsafe_privileged)>=1
    res := common.result({"original":json.marshal(unSafe_privileged), "Path": input.Path}, "KN-011")
}