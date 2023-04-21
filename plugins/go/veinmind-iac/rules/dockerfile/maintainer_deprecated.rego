package brightMirror.dockerfile

import data.common

risks[res]{
    inner:=input[_]
    inner.Cmd=="maintainer"
    res:=common.result(inner,"DF-022")
}