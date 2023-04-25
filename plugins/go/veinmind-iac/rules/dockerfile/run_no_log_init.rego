package brightMirror.dockerfile

import future.keywords.in
import data.common

risks[res] {
    inner := input[_]
	inner.Cmd == "run"
	some value in inner.Value
        contains(value,"useradd")
        not contains(value,"--no-log-init")
    res := common.result(inner, "DF-030")
}
