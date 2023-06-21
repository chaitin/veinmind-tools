package service

import (
	"strings"
	"unicode"
)

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
			continue
		}
		if checkUser(str) {
			risk = true
			riskContent += str + "\n"
			continue
		}
	}
	return risk, riskContent
}

// checkShell 检测反弹shell,下载执行
func checkShell(content string) bool {
	// 反弹shell
	if (strings.Contains(content, "bash") && (strings.Contains(content, "/dev/tcp/") || strings.Contains(content, "telnet ") || strings.Contains(content, "nc ") ||
		(strings.Contains(content, "exec ") && strings.Contains(content, "socket")) || strings.Contains(content, "curl ") ||
		strings.Contains(content, "wget ") || strings.Contains(content, "lynx ") || strings.Contains(content, "bash -i"))) ||
		strings.Contains(content, ".decode('base64')") || strings.Contains(content, "exec(base64.b64decode") {
		return true
	} else if strings.Contains(content, "/dev/tcp/") && (strings.Contains(content, "exec ") || strings.Contains(content, "ksh -c")) {
		return true
	} else if strings.Contains(content, "exec ") && (strings.Contains(content, "socket.") || strings.Contains(content, ".decode('base64')")) {
		return true
	}
	// 下载执行类
	if (strings.Contains(content, "wget ") || strings.Contains(content, "curl ")) &&
		(strings.Contains(content, " -O ") || strings.Contains(content, " -s ")) &&
		strings.Contains(content, " http") &&
		(strings.Contains(content, "php ") || strings.Contains(content, "perl") || strings.Contains(content, "python ") ||
			strings.Contains(content, "sh ") || strings.Contains(content, "bash ")) {
		return true
	}

	return false
}

// checkUser 检测用户修改行为
func checkUser(content string) bool {
	if strings.Contains(content, "useradd ") || strings.Contains(content, "usermod ") || strings.Contains(content, "userdel ") {
		return true
	}
	return false
}
