package service

import (
	"os"
	"strings"
	"unicode"
)

var malwareInfos []string

// analysisStrings 分析字符串是否为后门
func analysisStrings(fileContents string) (bool, string) {
	arr := strings.Split(fileContents, "\n")
	risk := false
	var riskContent string
	for _, str := range arr {
		str = strings.TrimLeftFunc(str, unicode.IsSpace)
		if len(str) == 0 || str[0] == '#' {
			continue
		}
		if checkShell(str) {
			risk = true
			riskContent += str + "\n"
		} else {
			for _, file := range strings.Fields(str) {
				if _, err := os.Stat(file); err != nil {
					continue
				}
				if info, _ := os.Stat(file); info.IsDir() {
					continue
				}
				// TODO： 如果在文件中引用其他文件，检测被引用文件内容是否为后门
			}
		}
	}
	return risk, riskContent
}

// checkShell 检测反弹shell
func checkShell(str string) bool {
	shellCommands := []string{
		"bash", "nc ", "curl ", "wget ", "ftp ", "telnet ",
		"ssh ", "scp ", "expect ", "php ", "jsp ",
		".php", ".jsp", "$(", "eval", "perl ", "python ",
		"ruby ", "lua ", "mysql ", "exec", "sh ", "/bin/sh ",
		"/usr/local/sbin", "powershell ", "select * from", "1;1", "$(echo",
		"source ", ". /etc/passwd", "update ", "delete ", "insert into",
		"set ", "grant ", "useradd", "usermod", "passwd ",
	}
	for _, command := range shellCommands {
		if strings.Contains(str, command) {
			return true
		}
	}
	return false
}
