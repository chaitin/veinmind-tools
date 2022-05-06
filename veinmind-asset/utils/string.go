package utils

func Limit(s string, num int) string {
	if len(s) > num {
		if num > 3 {
			return s[0:num-3] + "..."
		} else {
			return "..."
		}
	} else {
		return s
	}
}

func Repeat(s string, num int) string {
	res := ""
	for i := 0; i < num; i++ {
		res += s
	}
	return res
}
