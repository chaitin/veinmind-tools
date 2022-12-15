package git

import (
	api "github.com/go-git/go-git/v5"

	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func Clone(path string, url string, key string, insecure bool) error {
	log.Infof("start download %s at %s", url, path)
	opt := &api.CloneOptions{
		URL:             url,
		InsecureSkipTLS: insecure,
	}
	if key != "" {
		sshAuth, err := ssh.NewPublicKeysFromFile("git", key, "")
		if err != nil {
			return err
		}
		opt.Auth = sshAuth
	}
	_, err := api.PlainClone(path, false, opt)
	if err != nil {
		return err
	}
	return nil
}
