package service

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	api "github.com/chaitin/libveinmind/go"
)

func SudoCheck(fs api.FileSystem, fi os.FileInfo, filename string) (bool, error) {
	content, err := fs.Open("/etc/sudoers")
	if err != nil {
		return false, err
	}
	defer content.Close()
	scanner := bufio.NewScanner(content)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		}
		compile := regexp.MustCompile(SUDOREGEX)
		matches := compile.FindStringSubmatch(scanner.Text())
		if len(matches) == 3 {
			if matches[1] == "admin" || matches[1] == "sudo" || matches[1] == "root" { //sudo默认设置
				continue
			} else {
				sudoFile := matches[2]
				if strings.Contains(strings.ToLower(strings.TrimSpace(sudoFile)), filename) {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func init() {
	ImageCheckFuncMap["sudo"] = SudoCheck
	ContainerCheckFuncMap["sudo"] = SudoCheck
}
