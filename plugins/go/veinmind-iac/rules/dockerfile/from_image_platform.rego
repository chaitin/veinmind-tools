package brightMirror.dockerfile

import future.keywords.in
import data.common

risks[res] {
    inner := input[_]
	inner.Cmd == "from"
	some val in inner.Flags
        contains(val, "--platform")
	res := common.result(inner, "DF-004")
}