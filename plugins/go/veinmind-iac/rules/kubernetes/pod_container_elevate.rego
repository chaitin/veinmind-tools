package brightMirror.kubernetes

import data.common

risks[res] {
   containers := input.spec.containers
   securityContexts := containers[_].securityContext
   allowPrivilegeEscalations := securityContexts.allowPrivilegeEscalations
   allowPrivilegeEscalations == true
   res := common.result({"original": "UnSafeSettings:set allowPrivilegeEscalation=true", "Path": input.Path}, "KN-002")
}
