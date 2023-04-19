package brightMirror.dockerfile

import data.common

risks[res] {
    count(get_cmd) > 1
    obj := [cmd | cmd := get_cmd[_]; true]
    res := common.result(obj[1], "DF-016")
}

get_cmd[inner] {
	inner := input[_]
    inner.Cmd == "cmd"
}

