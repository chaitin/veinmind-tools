package harborhook

import (
	"testing"

	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/registry"
	"github.com/stretchr/testify/assert"
)

func TestParseHarborWebhookConfig(t *testing.T) {
	content := `[[policys]]
	action = "container_create"
	enabled_plugins = ["veinmind-weakpass"]
	plugin_params = ["veinmind-weakpass:scan.serviceName=ssh"]
	risk_level_filter = ["High"]
	block = true
	alert = true
	
	[log]
	report_log_path = "report.log"
	authz_log_path = "webhook.log"
	
	[docker_auth]
	registry = "index.docker.io"
	username = "huzai9527"
	password = "asdqwe123."
	
	[harbor_auth]
	registry = "10.9.33.98"
	username = "admin"
	password = "asdqwe123"	
	`
	config, err := ParseWebhookConfigFromReader(content)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, &WebhookConfig{
		Policys: []Policy{
			{
				Action:          "container_create",
				EnabledPlugins:  []string{"veinmind-weakpass"},
				PluginParams:    []string{"veinmind-weakpass:scan.serviceName=ssh"},
				RiskLevelFilter: []string{"High"},
				Block:           true,
				Alert:           true,
			},
		}, Log: Log{
			WebhookLogPath: "webhook.log",
			ReportLogPath:  "report.log",
		}, DockerAuth: registry.Auth{
			Registry: "index.docker.io",
			Username: "huzai9527",
			Password: "asdqwe123.",
		}, HarborAuth: registry.Auth{
			Registry: "10.9.33.98",
			Username: "admin",
			Password: "asdqwe123",
		},
	}, config)
}
