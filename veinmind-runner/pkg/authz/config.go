package authz

import (
	"errors"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/registry"
)

type Policy struct {
	Action          string   `toml:"action"`
	EnabledPlugins  []string `toml:"enabled_plugins"`
	PluginParams    []string `toml:"plugin_params"`
	RiskLevelFilter []string `toml:"risk_level_filter"`
	Block           bool     `toml:"block"`
	Alert           bool     `toml:"alert"`
}
type Log struct {
	ReportLogPath string `toml:"report_log_path"`
	AuthZLogPath  string `toml:"authz_log_path"`
}

type PolicyConfig struct {
	Policys    []Policy      `toml:"policys"`
	Log        Log           `toml:"log"`
	DockerAuth registry.Auth `toml:"docker_auth"`
}

func ParsePolicyConfig(path string) (*PolicyConfig, error) {
	policyConfig := &PolicyConfig{}
	if path == "" {
		return nil, errors.New("authz config path can't be empty")
	}
	_, err := toml.DecodeFile(path, policyConfig)
	if err != nil {
		return nil, err
	}

	return policyConfig, nil
}

// just for test
func ParsePolicyConfigFromReader(content string) (*PolicyConfig, error) {
	policyConfig := &PolicyConfig{}
	_, err := toml.Decode(content, policyConfig)
	if err != nil {
		return nil, err
	}

	return policyConfig, nil
}
func GetLogFile(output string) (*os.File, error) {
	if _, err := os.Stat(output); errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(output)
		if err != nil {
			return nil, err
		}
	}
	reportFile, err := os.OpenFile(output, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return reportFile, nil
}
func (t *PolicyConfig) PolicysMap() map[string]Policy {
	result := make(map[string]Policy)
	for _, t := range t.Policys {
		result[t.Action] = t
	}
	return result
}
