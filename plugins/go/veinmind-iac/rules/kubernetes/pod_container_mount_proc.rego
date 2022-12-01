package brightMirror.kubernetes

import data.common
import future.keywords.if
import future.keywords.in

unsafe_mount_proc[d]{
	inner := input.spec.volumes[i].hostPath
    some val in inner
    contains(val,"/proc")
    d:=input.spec.volumes[i].name
}
unSafe_mount_proc:={
    "UnSafeVolumnsName":unsafe_mount_proc
}
risks[res]{
    count(unsafe_mount_proc)>=1
    res := common.result({"original":json.marshal(unSafe_mount_proc), "Path": input.Path}, "KN-019")
}