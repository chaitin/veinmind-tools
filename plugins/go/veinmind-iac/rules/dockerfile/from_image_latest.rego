package brightMirror.dockerfile

import future.keywords.in
import data.common

# TODO far from perfect, improve

risks[res] {
    inner := input[_]
	inner.Cmd == "from"
	some val in inner.Value
        contains(val, ":latest")
        not contains(val, "scratch")
	res := common.result(inner, "DF-003")
}

risks[res] {
    inner := input[_]
	inner.Cmd == "from"
	some val in inner.Value
        not contains(val, ":")
        not contains(val, "scratch")
	res := common.result(inner, "DF-003")
}