package authz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAuthConfig(t *testing.T) {
	content := `[[policys]]
	action = "container_create"
	enabled_plugins = ["veinmind-weakpass"]
	plugin_params = ["veinmind-weakpass:scan.service=ssh"]
	risk_level_filter = ["High"]
	block = true
	alert = true
	[log]
	report_log_path = "report.log"
	authz_log_path = "authz.log"
	`
	config, err := ParsePolicyConfigFromReader(content)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, &PolicyConfig{Policys: []Policy{
		{
			Action:          "container_create",
			EnabledPlugins:  []string{"veinmind-weakpass"},
			PluginParams:    []string{"veinmind-weakpass:scan.service=ssh"},
			RiskLevelFilter: []string{"High"},
			Block:           true,
			Alert:           true,
		},
	}, Log: Log{
		AuthZLogPath:  "authz.log",
		ReportLogPath: "report.log",
	}}, config)
}
