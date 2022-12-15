package brightMirror.kubernetes

import data.common
import future.keywords.if
import future.keywords.in

risks[res]{
    inner := input.spec.volumes[i].hostPath
    some val in inner
    contains(val,"docker.sock")
    Name:=input.spec.volumes[i].name
    Names:=[Name]
    Hints:=["UnSafeVolumeName"]
    Combine:=array.concat(Hints,Names)
    res := common.result({"original":concat(":",Combine), "Path": input.Path}, "KN-016")
}