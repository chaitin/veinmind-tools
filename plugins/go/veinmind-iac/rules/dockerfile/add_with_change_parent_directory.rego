package brightMirror.dockerfile

import future.keywords.in
import data.common

risks[res]{
    inner:= input[_]
    inner.Cmd=="add"
    some val in inner.Value
        contains(val,"../")
    res := common.result(inner,"DF-015")
}