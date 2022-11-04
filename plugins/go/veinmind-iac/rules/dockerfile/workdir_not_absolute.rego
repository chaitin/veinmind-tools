package brightMirror.dockerfile

import future.keywords.in
import data.common

risks[res] {
    inner := input[_]
    inner.Cmd == "workdir"
    some dir in inner.Value
        not regex.match("^[\"']?(/[A-z0-9-_+]*)|([A-z0-9-_+]:\\\\.*)|(\\$[{}A-z0-9-_+].*)", dir)
    res := common.result(inner, "DF-006")
}