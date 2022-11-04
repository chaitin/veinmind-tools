package brightMirror.dockerfile

import data.common

risks[res] {
    inner := input[_]
	inner.Cmd == "expose"
	inner.Value[_] = ["22", "22/tcp", "22/udp"][_]
	res := common.result(inner, "DF-001")
}