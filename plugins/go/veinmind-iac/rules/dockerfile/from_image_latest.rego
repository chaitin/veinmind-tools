package brightMirror.dockerfile

import future.keywords.in
import data.common

# TODO far from perfect, improve

risks[res] {
    inner := input[_]
	inner.Cmd == "from"
	args := concat(" ", inner.Value)
    contains(args, ":latest")
    not contains(args, "scratch")
	res := common.result(inner, "DF-003")
}

risks[res] {
    inner := input[_]
	inner.Cmd == "from"
	args := concat(" ", inner.Value)
    not contains(args, ":")
    not contains(args, "scratch")
	res := common.result(inner, "DF-003")
}