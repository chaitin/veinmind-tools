package brightMirror.dockerfile

import data.common

dnf_install_regex := `(dnf install)|(dnf in)|(dnf reinstall)|(dnf rei)|(dnf install-n)|(dnf install-na)|(dnf install-nevra)`
dnf_regex = sprintf("(%s).*dnf.*clean.*all", [dnf_install_regex])

get_dnf[output] {
	run:=input[_]
	run.Cmd=="run"
	arg := run.Value[0]

	regex.match(dnf_install_regex, arg)

	not contains_clean_after_dnf(arg)
	output := run
}

risks[res] {
	output := get_dnf[_]
	res:=common.result(output, "DF-027")
}

contains_clean_after_dnf(cmd) {
	regex.match(dnf_regex, cmd)
}