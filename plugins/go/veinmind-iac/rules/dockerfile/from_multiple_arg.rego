package brightMirror.dockerfile

import data.common

risks[res]{
	inner_from:=input[_]
    inner_from.Cmd=="from"
    line:=inner_from.StartLine
    i<line
    PreArgs:=findPreArg[i]
    inner_argCheck:=input[_]
    inner_argCheck.StartLine>line
    contains(inner_argCheck.Original,PreArgs)
    res=common.result(inner_argCheck,"DF-021")
}

findPreArg[line]=ContainStrings {
	inner_argValue:=input[_]
    inner_argValue.Cmd=="arg"
    line:=inner_argValue.StartLine
    PreArgs:={PreArg|PreArg:=inner_argValue.Value}
    tmps:=PreArgs[_][_]
    tmp:=split(tmps,"=")
	ContainStrings=sprintf("${%s}",[tmp[0]])
}