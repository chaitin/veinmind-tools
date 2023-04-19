package brightMirror.dockerfile

import future.keywords.in
import data.common

risks[res] {
     inner := input[_]
     inner.Cmd == "run"
     some value in inner.Value
         find_update(value)
         not checkInvalid(value)
     res := common.result(inner, "DF-024")
}

find_update(command) {
	commands := regex.split(`\s*&&\s*`, command)
	array_split := split(commands[_], " ")
	len = count(array_split)
	update := {"update", "--update"}
	array_split[len - 1] == update[_]
}

checkInvalid(command) {
	command_list = [
		"upgrade",
		"install",
		"source-install",
		"reinstall",
		"groupinstall",
		"localinstall",
		"apk add",
	]

	update := indexof(command, "update")
	update != -1

	install := indexof(command, command_list[_])
	install != -1

	update < install
}