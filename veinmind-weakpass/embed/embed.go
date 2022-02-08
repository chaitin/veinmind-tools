package embed

import (
	"embed"
	common "github.com/chaitin/veinmind-tools/veinmind-weakpass/log"
	"io/ioutil"
	"os"
	"path"
)

//go:embed pass.dict
var EmbedFS embed.FS

func ExtractAll() {
	extract("pass.dict")
}

// extract
func extract(epath string) {
	// extract docker-compose config
	composeYamlBytes, err := EmbedFS.ReadFile(epath)
	if err != nil {
		common.Log.Fatal(err)
	}

	if _, err := os.Stat(path.Dir(epath)); os.IsNotExist(err) {
		err = os.Mkdir(path.Dir(epath), 0755)
		if err != nil {
			common.Log.Fatal(err)
		}
	}
	err = ioutil.WriteFile(epath, composeYamlBytes, 0755)
	if err != nil {
		common.Log.Fatal(err)
	}
}
