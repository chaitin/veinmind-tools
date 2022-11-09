package brightMirror.dockerfile

import data.common

risks[res] {
    inner := input[_]
    inner.Cmd == "add"
    args := concat(" ", inner.Value)
    not contains(args, ".tar")
    res := common.result(inner, "DF-010")
}