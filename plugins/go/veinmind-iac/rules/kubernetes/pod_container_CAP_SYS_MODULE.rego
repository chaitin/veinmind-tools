package brightMirror.kubernetes

import data.common
import future.keywords.in


unsafe_sys_module[d]{
	inner := input.spec.containers[i].securityContext.capabilities.add
    some val in inner
    upper(val) == "SYS_MODULE"
    d:=input.spec.containers[i].name
}
unSafe_sys_module:={
    "UnSafeContainersName":unsafe_sys_module
}
risks[res]{
    count(unsafe_sys_module)>=1
    res := common.result({"original":json.marshal(unSafe_sys_module), "Path": input.Path}, "KN-014")
}