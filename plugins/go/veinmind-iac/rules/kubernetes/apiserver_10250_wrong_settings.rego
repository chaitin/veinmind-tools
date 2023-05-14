package brightMirror.kubernetes

import data.common

risks[res]{
	input.authentication.anonymous.enabled==true
    input.authorization.mode=="AlwaysAllow"
    res := common.result({"original":"UnSafeSettings:`authentication.anonymous`,`authorization.mode`", "Path": input.Path}, "KN-007")
}