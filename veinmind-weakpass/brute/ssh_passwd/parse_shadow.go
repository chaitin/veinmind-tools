package ssh_passwd

import (
	"bufio"
	common "github.com/chaitin/veinmind-tools/veinmind-weakpass/log"
	"io"
	"strings"
)

type Shadow struct {
	LoginName                string
	EncryptedPassword        string
	LastPasswordChange       string
	MinimumPasswordAge       string
	MaximumPasswordAge       string
	PasswordWarningPeriod    string
	PasswordInactivityPeriod string
	AccountExpirationDate    string
	ReservedField            string
}

func ParseShadowFile(shadowFile io.Reader) (shadows []Shadow, err error) {
	scanner := bufio.NewScanner(shadowFile)
	for scanner.Scan() {
		userinfo := strings.Split(scanner.Text(), ":")
		if len(userinfo) != 9 {
			common.Log.Error("Shadow format error")
			continue
		}

		s := Shadow{}
		s.LoginName = userinfo[0]
		s.EncryptedPassword = userinfo[1]
		s.LastPasswordChange = userinfo[2]
		s.MinimumPasswordAge = userinfo[3]
		s.MaximumPasswordAge = userinfo[4]
		s.PasswordWarningPeriod = userinfo[5]
		s.PasswordInactivityPeriod = userinfo[6]
		s.AccountExpirationDate = userinfo[7]
		s.ReservedField = userinfo[8]

		shadows = append(shadows, s)
	}

	return shadows, nil
}
