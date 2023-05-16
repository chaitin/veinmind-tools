package brightMirror.kubernetes

import data.common

risks[res] {
	inner:= containers[i]
	key := sprintf("%s/%s", ["container.apparmor.security.beta.kubernetes.io", inner.name])
    annotations:=input[_].metadata.annotations
    annotations[key]!="runtime/default"
    res := common.result({"original": annotations[key],"Path": input[_].Path}, "KN-003")
}