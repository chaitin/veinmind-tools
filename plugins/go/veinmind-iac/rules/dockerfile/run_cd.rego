package brightMirror.dockerfile

import data.common

risks[res] {
    inner := input[_]
	inner.Cmd == "run"
    parts = regex.split(`\s*&&\s*`, inner.Value[_])
    startswith(parts[_], "cd ")
	res := common.result(inner, "DF-023")
}
