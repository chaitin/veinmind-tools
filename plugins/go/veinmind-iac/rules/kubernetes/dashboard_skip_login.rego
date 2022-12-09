package brightMirror.kubernetes

import data.common
import future.keywords.in
import future.keywords.contains
import future.keywords.if

risks[res]{
    contains(input.spec.containers[0].args[i],"enable-skip-login")
	res := common.result({"original":"UnSafeSettings:`spec.containers.args`", "Path": input.Path}, "KN-008")
}