package brightMirror.kubernetes

import data.common
import future.keywords.every
import future.keywords.in
import future.keywords.contains


risks[res]{
    containers[_].command[_]=="etcd"
	every val in containers[_].args{
    not contains(val,"--client-cert-auth=true")
    }
    res := common.result({"original":"UnSafeSettings:`spec.containers.command missing --client-cert-auth=true`", "Path": input[_].Path}, "KN-009")
}