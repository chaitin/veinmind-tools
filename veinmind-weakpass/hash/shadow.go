package hash

type Shadow struct {
	hash string
}

func (i *Shadow) ID() string {
	return "shadow"
}

func (i *Shadow) Plain() string {
	return i.hash
}

func (i *Shadow) Match(hash, guess string) (result string, flag bool) {
	var pwd Password
	if err := ParsePassword(&pwd, hash); err != nil {
		return "", false
	}
	if _,ok := pwd.Match([]string{guess});ok {
		return guess, true
	}
	return "", false
}
