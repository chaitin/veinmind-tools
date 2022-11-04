package brightMirror.kubernetes

import data.common

risks[res] {
    name := containers[_].name
	key := sprintf("%s/%s", ["container.apparmor.security.beta.kubernetes.io", name])
	val := annotations[i][key]
	val != "runtime/default"
    res := common.result({"original": val,"Path": input[i].Path}, "KN-003")
}