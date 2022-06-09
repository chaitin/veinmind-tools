package hash

type Plain struct {
}

func (i *Plain) ID() string {
	return "plain"
}

func (i *Plain) Match(hash, guess string) (flag bool, err error) {
	if hash == guess {
		return true, nil
	}
	return false, nil
}
