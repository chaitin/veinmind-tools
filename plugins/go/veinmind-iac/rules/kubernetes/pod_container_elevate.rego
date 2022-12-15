package brightMirror.kubernetes

import data.common

risks[res] {
   count(securityContexts) > 0
   count(allowPrivilegeEscalations) > 0
   allowPrivilegeEscalations[i] == true
   res := common.result({"original": allowPrivilegeEscalations[i], "Path": input[i].Path}, "KN-002")
}

risks[res] {
   count(securityContexts) > 0
   count(allowPrivilegeEscalations)  < 1
   res := common.result({"original":"UnSafeSettings:`unset allowPrivilegeEscalation=false`", "Path": input[i].Path}, "KN-002")
}