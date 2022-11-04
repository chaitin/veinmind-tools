package brightMirror.dockerfile

import data.common

get_user[inner] {
	inner := input[_]
    inner.Cmd == "user"
}

# no user command
risks[res] {
    count(get_user) < 1
	res := common.result({"Path": get_path[i]}, "DF-005")
}

risks[res] {
    users := [user | user := get_user[_]; true]
    equal(users[count(users) - 1].value[_], "root")
    res := common.result(users[count(users) - 1], "DF-005")
}