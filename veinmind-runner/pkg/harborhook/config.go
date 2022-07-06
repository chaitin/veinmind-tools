package harborhook

import (
	"errors"

	"github.com/BurntSushi/toml"
	"github.com/chaitin/veinmind-tools/veinmind-common/go/service/report"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/registry"
)

var (
	ToLevel = map[report.Level]string{
		report.Low:      "Low",
		report.Medium:   "Medium",
		report.High:     "High",
		report.Critical: "Critical",
		report.None:     "None",
	}

	FromLevel = map[string]report.Level{
		"Low":      report.Low,
		"Medium":   report.Medium,
		"High":     report.High,
		"Critical": report.Critical,
		"None":     report.None,
	}
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
	ReportLogPath  string `toml:"report_log_path"`
	WebhookLogPath string `toml:"authz_log_path"`
}
type WebhookConfig struct {
	Policys    []Policy      `toml:"policys"`
	Log        Log           `toml:"log"`
	DockerAuth registry.Auth `toml:"docker_auth"`
	HarborAuth registry.Auth `toml:"harbor_auth"`
}

func ParseWebHookConfig(path string) (*WebhookConfig, error) {
	WebhookConfig := &WebhookConfig{}
	if path == "" {
		return nil, errors.New("webhook config path can't be empty")
	}
	_, err := toml.DecodeFile(path, WebhookConfig)
	if err != nil {
		return nil, err
	}

	return WebhookConfig, nil
}

// just for test
func ParseWebhookConfigFromReader(content string) (*WebhookConfig, error) {
	WebhookConfig := &WebhookConfig{}
	_, err := toml.Decode(content, WebhookConfig)
	if err != nil {
		return nil, err
	}

	return WebhookConfig, nil
}

func (w *WebhookConfig) PolicysMap() map[string]Policy {
	result := make(map[string]Policy)
	for _, t := range w.Policys {
		result[t.Action] = t
	}
	return result
}
