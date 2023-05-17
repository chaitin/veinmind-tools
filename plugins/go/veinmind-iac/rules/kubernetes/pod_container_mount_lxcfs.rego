package brightMirror.kubernetes

import data.common
import future.keywords.if
import future.keywords.in

risks[res]{
    inner := volumes[_].hostPath
    some val in inner
    contains(val,"lxcfs")
    Name:=volumes[_].name
    Names:=[Name]
    Hints:=["UnSafeVolumeName"]
    Combine:=array.concat(Hints,Names)
    res := common.result({"original":concat(":",Combine), "Path": input[_].Path}, "KN-017")
}