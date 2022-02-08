package brute

import (
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/brute/ssh_passwd"
)

func SSHMatchPassword(encrypt string, guess string) (string, bool) {
	var pwd ssh_passwd.Password
	if err := ssh_passwd.ParsePassword(&pwd, encrypt); err != nil {
		return "", false
	}

	return pwd.Match([]string{guess})
}
