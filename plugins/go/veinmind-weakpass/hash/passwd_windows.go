package hash

import (
	"errors"
)

// PasswordMethod represents the entryption method.
type PasswordMethod uint8

// Password represents the data field inside.
type Password struct {
	// Method of current password.
	Method PasswordMethod

	// MethodString of the current password.
	MethodString string

	// Salt inside the password.
	Salt string

	// Hash hased value of the password.
	Hash string
}

func (pw *Password) Match(guesses []string) (string, bool) {
	panic("not implemented")
	return "", false
}

func ParsePassword(pass *Password, phrase string) error {
	return errors.New("not implemented")
}
