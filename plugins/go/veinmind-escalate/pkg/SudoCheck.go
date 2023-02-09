package pkg

import (
	"bufio"
	api "github.com/chaitin/libveinmind/go"
	"regexp"
	"strings"
)

func SudoFileCheck(fs api.FileSystem) error {
	UnsafeSudoFiles := []string{"wget", "find", "cat", "apt", "zip", "xxd", "time", "taskset", "git", "sed", "pip", "ed", "tmux", "scp", "perl", "bash", "less", "awk", "man", "vi", "vim", "env", "ftp", "all"}
	content, err2 := fs.Open("/etc/sudoers")
	if err2 != nil {
		return err2
	}
	defer func(err2 error) {
		if err2 == nil {
			content.Close()
		}
	}(err2)
	scanner := bufio.NewScanner(content)

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		}
		compile := regexp.MustCompile(SUDOREGEX)
		res := compile.FindStringSubmatch(scanner.Text())
		//fmt.Println(res)
		if len(res) == 3 {
			if res[1] == "admin" || res[1] == "sudo" || res[1] == "root" { //sudo默认设置
				continue
			} else { //其他用户
				sudoFile := res[2]
				for _, UnsafeSudoFile := range UnsafeSudoFiles {
					if strings.Contains(UnsafeSudoFile, strings.ToLower(strings.TrimSpace(sudoFile))) {
						AddResult(scanner.Text(), SUDOREASON, "UnSafeUser "+res[1])
					}
				}
			}
		}
	}
	return nil

}
