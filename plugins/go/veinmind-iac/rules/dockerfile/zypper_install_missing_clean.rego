package brightMirror.dockerfile

import data.common

install_regex := `(zypper in)|(zypper remove)|(zypper rm)|(zypper source-install)|(zypper si)|(zypper patch)|(zypper (-(-)?[a-zA-Z]+ *)*install)`

zypper_regex := sprintf("%s|(zypper clean)|(zypper cc)", [install_regex])

risks[res] {
	output := get_zypper[_]
	res :=common.result(output.cmd, "DF-029")
}

get_zypper[output] {
	run:=input[_]
    run.Cmd=="run"
	arg := run.Value[0]
	regex.match(install_regex, arg)
	not contains_zipper_clean(arg)
	output := {
		"arg": arg,
		"cmd": run
	}
}

contains_zipper_clean(cmd) {
	zypper_commands := regex.find_n(zypper_regex, cmd, -1)
	is_zypper_clean(zypper_commands[count(zypper_commands) - 1])
}

is_zypper_clean(cmd) {
	cmd == "zypper clean"
}

is_zypper_clean(cmd) {
	cmd == "zypper cc"
}