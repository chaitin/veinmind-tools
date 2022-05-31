package module

import (
	"bufio"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/ssh"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/dict"
	"io"
	"strings"
)

type Ssh struct {
	Module
}

func (this *Ssh) Init(conf Config) (err error) {
	this.Module.Init(conf)
	this.specialDict = dict.Sshdict
	return nil
}

func (this *Ssh) ParsePasswdInfo(shadowFile io.Reader) (shadows []PasswdInfo, err error) {
	scanner := bufio.NewScanner(shadowFile)
	for scanner.Scan() {
		userinfo := strings.Split(scanner.Text(), ":")
		if len(userinfo) != 9 {
			log.Warn("Shadow format error")
			continue
		}

		s := PasswdInfo{}
		s.Username = userinfo[0]
		s.Password = userinfo[1]
		shadows = append(shadows, s)
	}

	return shadows, nil
}

func (this *Ssh) MatchPasswd(encrypt string, guess string) bool {
	var pwd ssh.SSHPassword
	if err := ssh.ParseSSHPassword(&pwd, encrypt); err != nil {
		return false
	}

	return pwd.Match([]string{guess})
}

func init() {
	mod := &Ssh{}
	mod.name = "SSH"
	mod.passwdType = SSH
	mod.filePath = []string{"/etc/shadow"}
	Register(mod)
}
