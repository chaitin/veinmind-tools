package brightMirror.kubernetes

import data.common
import future.keywords.in


unsafe_cap_dac_read_search[d]{
	inner := input.spec.containers[i].securityContext.capabilities.add
    some val in inner
    upper(val) == "DAC_READ_SEARCH"
    d:=input.spec.containers[i].name
}
unSafe_cap_dac_read_search:={
    "UnSafeContainersName":unsafe_cap_dac_read_search
}
risks[res]{
    count(unsafe_cap_dac_read_search)>=1
    res := common.result({"original":json.marshal(unSafe_cap_dac_read_search), "Path": input.Path}, "KN-013")
}