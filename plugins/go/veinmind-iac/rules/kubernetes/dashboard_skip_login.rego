package brightMirror.kubernetes

import data.common
import future.keywords.contains

risks[res]{
    contains(containers[_].args[_],"enable-skip-login")
	res := common.result({"original":"UnSafeSettings:`spec.containers.args`", "Path": input[_].Path}, "KN-008")
}

risks[res]{
    contains(pods[_].spec.containers[_].args[_],"enable-skip-login")
	res := common.result({"original":"UnSafeSettings:`spec.containers.args`", "Path": input[_].Path}, "KN-008")
}