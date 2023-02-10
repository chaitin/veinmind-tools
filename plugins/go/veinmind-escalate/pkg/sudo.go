package pkg

import (
	"bufio"
	"github.com/chaitin/veinmind-common-go/service/report/event"
	"regexp"
	"strings"

	api "github.com/chaitin/libveinmind/go"
)

func SudoFileCheck(fs api.FileSystem) ([]*event.EscapeDetail, error) {
	var res = make([]*event.EscapeDetail, 0)
	UnsafeSudoFiles := []string{"wget", "find", "cat", "apt", "zip", "xxd", "time", "taskset", "git", "sed", "pip", "ed", "tmux", "scp", "perl", "bash", "less", "awk", "man", "vi", "vim", "env", "ftp", "all"}
	content, err := fs.Open("/etc/sudoers")
	if err != nil {
		return res, err
	}
	defer content.Close()
	scanner := bufio.NewScanner(content)

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		}
		compile := regexp.MustCompile(SUDOREGEX)
		matches := compile.FindStringSubmatch(scanner.Text())
		//fmt.Println(res)
		if len(matches) == 3 {
			if matches[1] == "admin" || matches[1] == "sudo" || matches[1] == "root" { //sudo默认设置
				continue
			} else { //其他用户
				sudoFile := matches[2]
				for _, UnsafeSudoFile := range UnsafeSudoFiles {
					if strings.Contains(UnsafeSudoFile, strings.ToLower(strings.TrimSpace(sudoFile))) {
						res = append(res, &event.EscapeDetail{
							Target: scanner.Text(),
							Reason: SUDOREASON,
							Detail: "UnSafeUser " + matches[1],
						})
					}
				}
			}
		}
	}
	return res, nil

}

func init() {
	ContainerCheckList = append(ContainerCheckList, SudoFileCheck)
	ImageCheckList = append(ImageCheckList, SudoFileCheck)
}
