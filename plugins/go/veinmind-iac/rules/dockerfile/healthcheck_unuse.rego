package brightMirror.dockerfile

import data.common

get_health[inner] {
	inner := input[_]
    inner.Cmd == "healthcheck"
}

get_path[path] {
    inner := input[_]
    path = inner.Path
}

risks[res] {
    count(get_health) < 1
    res = common.result({"Path": get_path[i]}, "DF-007")
}

risks[res] {
    count(get_health) == 1
    healthes := [health | health := get_health[_]; true]
    healthes[0].Value == null
    res := common.result(healthes[0], "DF-007")
}