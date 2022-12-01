package brightMirror.kubernetes

import data.common
import future.keywords.in


unsafe_sys_ptrace[d]{
    input.spec.hostPID==true
	inner := input.spec.containers[i].securityContext.capabilities.add
    some val in inner
    upper(val) == "SYS_PTRACE"
    d:=input.spec.containers[i].name
}
unSafe_sys_ptrace:={
    "UnSafeContainersName":unsafe_sys_ptrace
}
risks[res]{
    count(unsafe_sys_ptrace)>=1
    res := common.result({"original":json.marshal(unSafe_sys_ptrace), "Path": input.Path}, "KN-020")
}