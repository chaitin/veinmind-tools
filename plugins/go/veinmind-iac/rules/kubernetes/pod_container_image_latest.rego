package brightMirror.kubernetes

import data.common

risks[res] {
	image := input.spec.containers[_].image
	contains(image, "latest")
    res := common.result({"original":input.spec.containers[_].image, "Path": input.Path}, "KN-001")
}

risks[res] {
    image := input.spec.containers[_].image
    not contains(image, ":")
    not equal(image, "scratch")
    res := common.result({"original":input.spec.containers[_].image, "Path": input.Path}, "KN-001")
}