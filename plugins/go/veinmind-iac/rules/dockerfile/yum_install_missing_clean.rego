package brightMirror.dockerfile

import data.common

yum_install_regex := `yum (-[a-zA-Z]+ *)*install`
yum_regex = sprintf("(%s).*yum.*clean.*all", [yum_install_regex])

get_yum[output] {
	run:=input[_]
	run.Cmd=="run"
	arg := run.Value[0]

	regex.match(yum_install_regex, arg)

	not contains_clean_after_yum(arg)
	output := run
}

risks[res] {
	output := get_yum[_]
	res:=common.result(output, "DF-028")
}

contains_clean_after_yum(cmd) {
	regex.match(yum_regex, cmd)
}