package brightMirror.dockerfile

import future.keywords.in
import data.common

short_flags := `(-([a-xzA-XZ])*y([a-xzA-XZ])*)`
long_flags := `(--yes)|(--assume-yes)`
optional_not_related_flags := `\s*(-(-)?[a-zA-Z]+\s*)*`
combined_flags := sprintf(`%s(%s|%s)%s`, [optional_not_related_flags, short_flags, long_flags, optional_not_related_flags])

risks[res]{
    output := get_apt_get[_]
	res:= common.result(output.cmd, "DF-025")
}

is_apt_get(command) {
	regex.match("(apt-get|yum|apt) (-(-)?[a-zA-Z]+ *)*install(-(-)?[a-zA-Z]+ *)*", command)
}

get_apt_get[output] {
	run:=input[_]
    run.Cmd=="run"
	count(run.Value) > 1
	arg := concat(" ", run.Value)
	is_apt_get(arg)
	not includes_assume_yes(arg)
	output := {
		"arg": arg,
		"cmd": run,
	}
}

get_apt_get[output] {
	run:=input[_]
    run.Cmd=="run"
	count(run.Value) == 1
	arg := run.Value[0]
	is_apt_get(arg)
	not includes_assume_yes(arg)
	output := {
		"arg": arg,
		"cmd": run,
	}
}

# checking json array
get_apt_get[output] {
	run:=input[_]
    run.Cmd=="run"
	count(run.Value) > 1
	arg := concat(" ", run.Value)
	is_apt_get(arg)
	not includes_assume_yes(arg)
	output := {
		"arg": arg,
		"cmd": run,
	}
}


# flags before command
includes_assume_yes(command) {
	install_regexp := sprintf(`(apt-get|yum|apt)%sinstall`, [combined_flags])
	regex.match(install_regexp, command)
}

# flags behind command
includes_assume_yes(command) {
	install_regexp := sprintf(`(apt-get|yum|apt)%sinstall%s`, [optional_not_related_flags, combined_flags])
	regex.match(install_regexp, command)
}