package ssh_passwd

import (
	"testing"
)

func TestParsePassword(t *testing.T) {
	// Encrypted shadow password using MD5.
	encrypt := "$1$Bg1H/4mz$X89TqH7tpi9dX1B9j5YsF."

	// Parse the password from the encryption.
	var pwd Password
	if err := ParsePassword(&pwd, encrypt); err != nil {
		t.Errorf("test: parsing password: %s", err)
		return
	}

	// Match the password with its original text.
	guesses := []string{"1234", "3456", "123"}
	if ma, m := pwd.Match(guesses); !m {
		t.Errorf("test: expecting match")
		return
	} else if ma != "123" {
		t.Errorf("test: wrong guessing: %s", ma)
		return
	}
}
