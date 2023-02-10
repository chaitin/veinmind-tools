package config

import (
	"os"
	"strings"

	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/joho/godotenv"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/embed"
)

func init() {
	// 加载配置
	f, err := embed.Open("scripts/.env")
	if err != nil {
		log.Fatal(err)
	}

	env, err := godotenv.Parse(f)
	if err != nil {
		log.Fatal(err)
	}

	currentEnv := map[string]bool{}
	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}

	for k, v := range env {
		if !currentEnv[k] {
			os.Setenv(k, v)
		}
	}
}
