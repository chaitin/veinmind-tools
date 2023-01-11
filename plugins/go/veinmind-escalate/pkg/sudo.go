package pkg

import (
	"bufio"
	api "github.com/chaitin/libveinmind/go"
	"strings"
)

const (
	SUDOREASON = "This file is granted sudo privileges and can be used for escalating"
)

func SudoFileCheck(fs api.FileSystem) error {
	UnsafeSudoFiles := []string{"wget", "find", "cat", "apt", "zip", "xxd", "time", "taskset", "git", "sed", "pip", "ed", "tmux", "scp", "perl", "bash", "less", "awk", "man", "vi", "vim", "env", "ftp", "ALL"}
	content, err := fs.Open("/etc/sudoers")
	defer content.Close()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(content)

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		}
		slice := strings.Split(scanner.Text(), " ")
		if slice[0] != "root" {
			sudoFile := slice[len(slice)-1]
			for _, UnsafeSudoFile := range UnsafeSudoFiles {
				if strings.Contains(UnsafeSudoFile, sudoFile) {
					AddResult(sudoFile, SUDOREASON, "UnSafeUser "+slice[0])
				}
			}

		}
	}
	return nil

}
