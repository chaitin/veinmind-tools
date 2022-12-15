package brightMirror.kubernetes

import data.common
import future.keywords.in
import future.keywords.contains
import future.keywords.if

risks[res]{
	input.authentication.anonymous.enabled==true
    input.authorization.mode=="AlwaysAllow"
    res := common.result({"original":"UnSafeSettings:`authentication.anonymous`,`authorization.mode`", "Path": input.Path}, "KN-007")
}