package service

import (
	"os"
	"strings"
)

var malwareInfos []string

// analysisStrings 分析字符串是否为后门
func analysisStrings(contents string) (bool, string) {
	content := strings.ReplaceAll(contents, "\n", "")
	if checkShell(content) {
		return true, content
	} else {
		for _, file := range strings.Fields(content) {
			if _, err := os.Stat(file); err != nil {
				continue
			}
			if info, _ := os.Stat(file); info.IsDir() {
				continue
			}
			// TODO： 如果在文件中引用其他文件，检测被引用文件内容是否为后门
		}
		return false, ""
	}
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
