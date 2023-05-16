package brightMirror.kubernetes

import data.common

risks[res]{
	input[_].authentication.anonymous.enabled==true
    input[_].authorization.mode=="AlwaysAllow"
    res := common.result({"original":"UnSafeSettings:`authentication.anonymous`,`authorization.mode`", "Path": input[_].Path}, "KN-007")
}