package brightMirror.kubernetes

import data.common

risks[res] {
   allowPrivilegeEscalations[_] == true
   res := common.result({"original": "UnSafeSettings:set allowPrivilegeEscalation=true", "Path": input[_].Path}, "KN-002")
}
