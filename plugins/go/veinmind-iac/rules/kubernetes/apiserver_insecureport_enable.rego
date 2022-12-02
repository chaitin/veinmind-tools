package brightMirror.kubernetes

import data.common
import future.keywords.in
import future.keywords.contains
import future.keywords.if

risks[res]{
    version:=input.spec.containers[i].image
    contains(version,"v1.1")
	inner:=input.spec.containers[i].command
    some val in inner
        contains(val,"insecure-port")
        code:=val
   res := common.result({"original":code, "Path": input.Path}, "KN-005")
}