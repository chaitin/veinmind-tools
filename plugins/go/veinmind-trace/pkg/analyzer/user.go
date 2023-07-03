package analyzer

import (
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/passwd"
)

func checkUsers(fs api.FileSystem) error {
	entries, err := passwd.ParseFilesystemPasswd(fs)
	entryMap := make(map[string]string, 0)
	if err != nil {
		return err
	}
	for _, e := range entries {
		// 1. check uid=0 but not root user
		if e.Uid == "0" && e.Username != "root" {
			// todo: may trace user, warn
		}
		// 2. check gid=0 but not root user
		if e.Gid == "0" && e.Username != "root" {
			// todo: may trace user, warn
		}
		// 3. check same uid user
		if _, ok := entryMap[e.Uid]; ok && e.Username != entryMap[e.Uid] {
			// todo: may trace add user, warn
		} else {
			entryMap[e.Uid] = e.Username
		}

	}
	return nil
}
