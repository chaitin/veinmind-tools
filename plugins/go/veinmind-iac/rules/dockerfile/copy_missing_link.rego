package brightMirror.dockerfile

import future.keywords.every
import data.common

risks[res] {
    inner := input[_]
	inner.Cmd == "copy"
	every flag in inner.Flags{
	    flag !="--link"
	}
	res := common.result(inner, "DF-020")
}