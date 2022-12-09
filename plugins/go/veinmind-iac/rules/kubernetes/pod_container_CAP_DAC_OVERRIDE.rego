package brightMirror.kubernetes

import data.common
import future.keywords.in

risks[res]{
        inner := input.spec.containers[i].securityContext.capabilities.add
        some val in inner
        upper(val) == "DAC_OVERRIDE"
        Name:=input.spec.containers[i].name
        Hints=["UnsafeContainers"]
        Names=[Name]
        Combine:=array.concat(Hints,Names)
        res := common.result({"original":concat(":",Combine), "Path": input.Path}, "KN-015")

}