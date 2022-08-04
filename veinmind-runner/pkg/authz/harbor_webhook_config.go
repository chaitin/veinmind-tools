package authz

import (
	"errors"

	"github.com/BurntSushi/toml"
	"github.com/chaitin/veinmind-common-go/pkg/auth"
)

var (
	defaultWebHookServer = WebhookServer{Port: 8080}
)

type MailConf struct {
	Host       string   `toml:"host"`
	Port       int      `toml:"port"`
	Name       string   `toml:"username"`
	Password   string   `toml:"password"`
	Subscriber []string `toml:"subscriber"`
}

type HarborPolicy struct {
	Policy
	SendMail bool `toml:"send_mail"`
}

type WebhookServer struct {
	Port          int    `toml:"port"`
	Authorization string `toml:"authorization"`
}

type HarborWebhookConfig struct {
	Log           Log            `toml:"log"`
	WebhookServer WebhookServer  `toml:"webhook_server"`
	DockerAuth    auth.Auth      `toml:"docker_auth"`
	Policies      []HarborPolicy `toml:"policies"`
	MailConf      MailConf       `toml:"mail_conf"`
}

func NewHarborWebhookConfig(paths ...string) (*HarborWebhookConfig, error) {
	if len(paths) < 1 {
		return nil, errors.New("config path can't be empty")
	}

	path := defaultConfigPath
	if paths[0] != "" {
		path = paths[0]
	}

	result := &HarborWebhookConfig{}
	_, err := toml.DecodeFile(path, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
