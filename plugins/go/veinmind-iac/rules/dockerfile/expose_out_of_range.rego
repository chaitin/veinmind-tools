package brightMirror.dockerfile

import data.common

risks[res] {
    inner := input[_]
	inner.Cmd == "expose"
	to_number(inner.Value[_]) > 65535
	res := common.result(inner, "DF-002")
}


risks[res] {
    inner := input[_]
	inner.Cmd == "expose"
	to_number(inner.Value[_]) <= 0
	res := common.result(inner, "DF-002")
}