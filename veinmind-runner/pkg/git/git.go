package git

import (
	"github.com/chaitin/libveinmind/go/plugin/log"
	api "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStr(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}

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
