package brightMirror.dockerfile

import future.keywords.in
import data.common

risks[res] {
    inner := input[_]
    inner.Cmd == "copy"
    some val in inner.Flags
        contains(val, "--chown")
    res := common.result(inner, "DF-009")
}