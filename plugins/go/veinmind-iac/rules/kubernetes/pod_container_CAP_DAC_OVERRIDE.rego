package brightMirror.kubernetes

import data.common
import future.keywords.in


unsafe_dac_override[d]{
	inner := input.spec.containers[i].securityContext.capabilities.add
    some val in inner
    upper(val) == "DAC_OVERRIDE"
    d:=input.spec.containers[i].name
}
unSafe_dac_override:={
    "UnSafeContainersName":unsafe_dac_override
}
risks[res]{
    count(unsafe_dac_override)>=1
    res := common.result({"original":json.marshal(unSafe_dac_override), "Path": input.Path}, "KN-015")
}