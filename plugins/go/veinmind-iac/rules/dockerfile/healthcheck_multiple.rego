package brightMirror.dockerfile

import data.common

get_health[inner] {
	inner := input[_]
    inner.Cmd == "healthcheck"
}

risks[res] {
    count(get_health) > 1
    obj := [health | health := get_health[_]; true]
    res := common.result(obj[1], "DF-008")
}
