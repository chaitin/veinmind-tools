package hash

type Shadow struct {
}

func (i *Shadow) ID() string {
	return "shadow"
}

func (i *Shadow) Match(hash, guess string) (flag bool, err error) {
	var pwd Password
	if err := ParsePassword(&pwd, hash); err != nil {

		return false, err
	}
	if _, ok := pwd.Match([]string{guess}); ok {
		return true, nil
	}
	return false, nil
}
