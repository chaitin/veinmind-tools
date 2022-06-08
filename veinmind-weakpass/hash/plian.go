package hash

type Plain struct {
	plain string
}

func (i *Plain) ID() string {
	return "plain"
}

func (i *Plain) Plain() (plain string) {
	return i.plain
}

func (i *Plain) Match(hash, guess string) (result string, flag bool) {
	if hash == guess {
		return hash, true
	}
	return "", false
}
