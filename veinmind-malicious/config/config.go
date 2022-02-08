package config

import (
	"github.com/chaitin/veinmind-tools/veinmind-malicious/embed"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/sdk/common"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

func init() {
	// 加载配置
	f, err := embed.Open("scripts/.env")
	if err != nil {
		common.Log.Fatal(err)
	}

	env, err := godotenv.Parse(f)
	if err != nil {
		common.Log.Fatal(err)
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
