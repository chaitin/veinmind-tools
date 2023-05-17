package brightMirror.kubernetes

import data.common

risks[res] {
	image := containers[_].image
	contains(image, "latest")
    res := common.result({"original":containers[_].image, "Path": input[_].Path}, "KN-001")
}

risks[res] {
    image := containers[_].image
    not contains(image, ":")
    not equal(image, "scratch")
    res := common.result({"original":containers[_].image, "Path": input[_].Path}, "KN-001")
}