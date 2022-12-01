package brightMirror.kubernetes

import future.keywords.every
import data.common
import future.keywords.in
import future.keywords.contains
import future.keywords.if


risks[res]{
    input.spec.containers[i].command[i]=="etcd"
	every val in input.spec.containers[i].command{
    not contains(val,"--peer-client-cert-auth=true")
    }
    res := common.result({"original":"missing --peer-client-cert-auth=true", "Path": input.Path}, "KN-010")
}