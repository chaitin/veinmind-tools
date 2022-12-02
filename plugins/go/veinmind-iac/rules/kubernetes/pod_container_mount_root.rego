package brightMirror.kubernetes

import data.common
import future.keywords.if
import future.keywords.in

unsafe_mount_root[d]{
	inner := input.spec.volumes[i].hostPath
    some val in inner
    val=="/"
    d:=input.spec.volumes[i].name
}
unSafe_mount_root:={
    "UnSafeVolumnsName":unsafe_mount_root
}
risks[res]{
    count(unsafe_mount_root)>=1
    res := common.result({"original":json.marshal(unSafe_mount_root), "Path": input.Path}, "KN-018")
}