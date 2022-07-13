package authz

import (
	"errors"
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
	AuthZLogPath  string `toml:"auth_log_path"`
	PluginLogPath string `toml:"plugin_log_path"`
}

type Listener struct {
	ListenAddr string `toml:"listener_addr"`
}

type DockerPluginConfig struct {
	Log        Log           `toml:"log"`
	Listener   Listener      `toml:"listener"`
	DockerAuth registry.Auth `toml:"docker_auth"`
	Policies   []Policy      `toml:"policies"`
}

const (
	defaultConfigPath = "config.toml"
)

func NewDockerPluginConfig(paths ...string) (*DockerPluginConfig, error) {
	if len(paths) < 1 {
		return nil, errors.New("config path can't be empty")
	}

	path := defaultConfigPath
	if paths[0] != "" {
		path = paths[0]
	}

	path = defaultConfigPath

	result := &DockerPluginConfig{}
	_, err := toml.DecodeFile(path, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
