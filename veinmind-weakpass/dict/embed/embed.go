package embed

import (
	"embed"
	"io/ioutil"
	"os"
	"path"
)

//go:embed pass.dict redis.dict tomcat.dict ssh.dict
var EmbedFS embed.FS

func ExtractAll() {
	extract("pass.dict")
	extract("tomcat.dict")
	extract("redis.dict")
	extract("ssh.dict")
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
	err = ioutil.WriteFile(epath, composeYamlBytes, 0755)
	if err != nil {
		return err
	}
	return nil
}
