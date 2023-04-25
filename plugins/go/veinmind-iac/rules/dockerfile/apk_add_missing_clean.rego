package brightMirror.dockerfile

import data.common

get_apk[output] {
	run := input[_]
    run.Cmd=="run"
	arg := run.Value[0]
	regex.match("apk (-[a-zA-Z]+\\s*)*add", arg)
	not contains_no_cache(arg)
	output=run
}

risks[res] {
	output := get_apk[_]
	res:=common.result(output, "DF-026")
}

contains_no_cache(cmd) {
	split(cmd, " ")[_] == "--no-cache"
}