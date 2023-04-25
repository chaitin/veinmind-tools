package brightMirror.dockerfile

import future.keywords.in
import data.common

risks[res] {
    count(get_entrypoint) > 1
    obj := [entrypoint | entrypoint := get_entrypoint[_]; true]
    res := common.result(obj[1], "DF-017")
}

get_entrypoint[inner] {
	inner := input[_]
    inner.Cmd == "entrypoint"
}