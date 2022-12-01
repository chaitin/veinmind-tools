package brightMirror.kubernetes

import data.common
import future.keywords.in
import future.keywords.contains
import future.keywords.if

risks[res]{
    input.spec.template.spec.containers[0].args[i]=="--enable-skip-login"
	res := common.result({"original":input.spec.template.spec.containers[0].args[i], "Path": input.Path}, "KN-008")
}