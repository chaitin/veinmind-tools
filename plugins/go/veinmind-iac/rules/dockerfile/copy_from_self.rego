package brightMirror.dockerfile

import data.common

find_alias_from_copy[out]{
    inner:=input[_]
    inner.Cmd=="copy"
    flags:=inner.Flags[_]
    contains(flags,"from")
    parts := split(flags, "=")
    is_equal(inner.Stage,parts[1])
    out:=inner
}

is_equal(stage,alias)=allow{
	inner:=input[_]
    inner.Stage==stage
    inner.Cmd="from"
    val:=inner.Value
    val[i]="as"
    current_alias:=val[i+1]
    current_alias==alias
    allow=true
}

risks[res]{
    out:=find_alias_from_copy[_]
	res:= common.result(out, "DF-019")
}
