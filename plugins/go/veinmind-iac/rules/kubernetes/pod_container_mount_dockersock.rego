package brightMirror.kubernetes

import data.common
import future.keywords.if
import future.keywords.in

unsafe_mount_dockersock[d]{
	inner := input.spec.volumes[i].hostPath
    some val in inner
    contains(val,"docker.sock")
    d:=input.spec.volumes[i].name
}
unSafe_mount_dockersock:={
    "UnSafeVolumnsName":unsafe_mount_dockersock
}
risks[res]{
    count(unsafe_mount_dockersock)>=1
    res := common.result({"original":json.marshal(unSafe_mount_dockersock), "Path": input.Path}, "KN-016")
}