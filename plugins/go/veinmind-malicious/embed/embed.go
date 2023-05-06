package embed

import (
	"embed"
	"io/fs"
	"os"
	"path"

	"github.com/chaitin/libveinmind/go/plugin/log"
)

//go:embed scripts/.env
//go:embed template/template.html
//go:embed scripts/docker-compose.yml
var EmbedFile embed.FS

func Open(name string) (fs.File, error) {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return EmbedFile.Open(name)
	} else {
		return os.Open(name)
	}
}

func ReadFile(name string) ([]byte, error) {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return EmbedFile.ReadFile(name)
	} else {
		return os.ReadFile(name)
	}
}

func ExtractAll() {
	extract("scripts/docker-compose.yml")
	extract("scripts/.env")
}

// extract
func extract(epath string) {
	// extract docker-compose config
	composeYamlBytes, err := EmbedFile.ReadFile(epath)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(path.Dir(epath)); os.IsNotExist(err) {
		err = os.Mkdir(path.Dir(epath), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = os.WriteFile(epath, composeYamlBytes, 0755)
	if err != nil {
		log.Fatal(err)
	}
}
