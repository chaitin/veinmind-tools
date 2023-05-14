package brightMirror.kubernetes

import data.common

risks[res] {
	inner:= input.spec.containers[i]
	key := sprintf("%s/%s", ["container.apparmor.security.beta.kubernetes.io", inner.name])
    annotations:=input.metadata.annotations
    annotations[key]!="runtime/default"
    res := common.result({"original": annotations[key],"Path": input.Path}, "KN-003")
}