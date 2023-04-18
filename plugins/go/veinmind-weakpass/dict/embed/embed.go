package embed

import (
	"embed"
	"os"
	"path"

	"github.com/chaitin/libveinmind/go/plugin/log"
)

//go:embed pass.dict redis.dict tomcat.dict ssh.dict ftp.dict
var EmbedFS embed.FS

func ExtractAll() {
	err := extract("pass.dict")
	if err != nil {
		log.Error(err)
	}
	err = extract("tomcat.dict")
	if err != nil {
		log.Error(err)
	}
	err = extract("redis.dict")
	if err != nil {
		log.Error(err)
	}
	err = extract("ssh.dict")
	if err != nil {
		log.Error(err)
	}
	err = extract("ftp.dict")
	if err != nil {
		log.Error(err)
	}
}

// extract
func extract(epath string) error {
	// extract docker-compose config
	composeYamlBytes, err := EmbedFS.ReadFile(epath)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path.Dir(epath)); os.IsNotExist(err) {
		err = os.Mkdir(path.Dir(epath), 0755)
		if err != nil {
			return err
		}
	}
	err = os.WriteFile(epath, composeYamlBytes, 0755)
	if err != nil {
		return err
	}
	return nil
}
