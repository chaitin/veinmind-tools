package brightMirror.dockerfile

import future.keywords.in
import data.common

get_aliased_name[output] {
	inner:=input[_]
	inner.Cmd=="from"
	value:=inner.Value
  	value[i]=="as"
	output={
    	"cmd":inner,
    	"startLine":inner.StartLine,
        "alias":value[i+1]
    }
}

checkDuplicate[output]{
    alias1:=get_aliased_name[_]
    alias2:=get_aliased_name[_]
    alias1.startLine!=alias2.startLine
    alias1.alias==alias2.alias
    output:=alias1.cmd
}

risks[res]{
	count(checkDuplicate)>1
    finalResults:=[finalResult|finalResult:=checkDuplicate[_]]
    index:=count(finalResults)-1
    res:=common.result(finalResults[index],"DF-018")
}