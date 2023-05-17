package brightMirror.kubernetes

import data.common
import future.keywords.in
import future.keywords.contains
import future.keywords.if


risks[res]{
	    containers[_].command[_]=="kube-apiserver"
	    version:=containers[_].image
	    contains(version,"v1.1")
	    not contains(version,"v1.19")
		inner:=containers[_].args
	    some val in inner
	        contains(val,"insecure-port")
	        not contains(val,"insecure-port=0")
	    res := common.result({"original":"UnSafeSettings:`spec.containers.args", "Path": input[_].Path}, "KN-005")
}

