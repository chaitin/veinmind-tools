package brightMirror.kubernetes

import data.common
import future.keywords.every

default allow=false

allowValuesSeLinuxOptionsType:=["container_t","container_init_t","container_kvm_t"]

risks[res]{
    type := securityContexts[_].seLinuxOptions.type
    every val in allowValuesSeLinuxOptionsType {
        val != type
    }
    res:= common.result({"original":"UnSafeSettings:`spec.containers.securityContext.seLinuxOptions.type`", "Path": input[_].Path},"KN-004")
}