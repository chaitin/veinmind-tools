package brightMirror.kubernetes

import data.common
import future.keywords.in

risks[res] {
	inner := input.spec.containers[i].securityContext.capabilities.add
	some val in inner
	val == "SYS_ADMIN"
	res :=common.result({"original":concat(" ",input.spec.containers[i].securityContext.capabilities), "Path": input.Path}, "KN-012")
}

unsafe_sysadmin[d]{
	inner := input.spec.containers[i].securityContext.capabilities.add
    some val in inner
    upper(val) == "SYS_ADMIN"
    d:=input.spec.containers[i].name
}
unSafe_sys_admin:={
    "UnSafeContainersName":unsafe_sysadmin
}
risks[res]{
    count(unsafe_sysadmin)>=1
    res := common.result({"original":json.marshal(unSafe_sys_admin), "Path": input.Path}, "KN-012")
}