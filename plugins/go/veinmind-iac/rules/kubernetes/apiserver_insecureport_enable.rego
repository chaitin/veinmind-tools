package brightMirror.kubernetes

import data.common
import future.keywords.in
import future.keywords.contains
import future.keywords.if

risks[res]{
    input.spec.containers[i].command[i]=="kube-apiserver"
    version:=input.spec.containers[i].image
    contains(version,"v1.1")
	inner:=input.spec.containers[i].command
    some val in inner
        contains(val,"insecure-port")
        not contains(val,"insecure-port=0")
        code:=val
   res := common.result({"original":"UnSafeSettings:`spec.containers.command`", "Path": input.Path}, "KN-005")
}
