package brightMirror.dockerfile

import future.keywords.in
import data.common

risks[res] {
    inner := input[_]
    some val in inner
        contains(val.Original, "BEGIN EC PRIVATE KEY")
    res := common.result(inner, "DF-009")
}