package module

type PasswordType uint8

const (
	SSH PasswordType = iota
	TOMCAT
	REDIS
)

type PasswdInfo struct {
	Username string
	Password string
	Filepath string
}

type BruteOption struct {
	Guess      string
	Passwdinfo PasswdInfo
}
