package brightMirror.dockerfile

import future.keywords.in
import data.common

risks[res] {
    inner := input[_]
    inner.Cmd == "run"
    some val in inner.Value
        parts = split(val, "&&")
    	instruction := parts[_]
        regex.match(`^\s*sudo`, instruction)
    res := common.result(inner, "DF-014")
}