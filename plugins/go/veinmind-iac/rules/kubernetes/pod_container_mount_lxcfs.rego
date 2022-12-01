package brightMirror.kubernetes

import data.common
import future.keywords.if
import future.keywords.in

unsafe_mount_lxcfs[d]{
	inner := input.spec.volumes[i].hostPath
    some val in inner
    contains(val,"lxcfs")
    d:=input.spec.volumes[i].name
}
unSafe_mount_lxcfs:={
    "UnSafeVolumnsName":unsafe_mount_lxcfs
}
risks[res]{
    count(unsafe_mount_lxcfs)>=1
    res := common.result({"original":json.marshal(unSafe_mount_lxcfs), "Path": input.Path}, "KN-017")
}