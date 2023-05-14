package brightMirror.kubernetes

import data.common
import future.keywords.contains

risks[res]{
    contains(input.spec.containers[_].args[i],"enable-skip-login")
	res := common.result({"original":"UnSafeSettings:`spec.containers.args`", "Path": input.Path}, "KN-008")
}

risks[res]{
    contains(input.spec.template.spec.containers[_].args[i],"enable-skip-login")
	res := common.result({"original":"UnSafeSettings:`spec.containers.args`", "Path": input.Path}, "KN-008")
}