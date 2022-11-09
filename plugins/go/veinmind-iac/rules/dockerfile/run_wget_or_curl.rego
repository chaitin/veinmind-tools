package brightMirror.dockerfile

import future.keywords.in
import data.common

get_curl[inner] {
    inner := input[_]
    inner.Cmd == "run"
    some val in inner.Value
        contains(val, "curl")
}

get_wget[inner] {
    inner := input[_]
    inner.Cmd == "run"
    some val in inner.Value
        contains(val, "wget")
}

risks[res] {
    count(get_curl) >= 1
    count(get_wget) >= 1
    command := [user | user := get_curl[_]; true]
    res := common.result(command[0], "DF-013")
}