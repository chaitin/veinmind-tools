package registry

import (
	"errors"
	"github.com/BurntSushi/toml"
	"net/url"
	"strings"
)

type Auth struct {
	Registry string `toml:"registry"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type AuthConfig struct {
	Auths []Auth `toml:"auths"`
}

func parseAuthConfig(path string) (*AuthConfig, error) {
	authConfig := &AuthConfig{}
	_, err := toml.DecodeFile(path, authConfig)
	if err != nil {
		return nil, err
	}

	return authConfig, nil
}

func parseAuthConfigFromString(content string) (*AuthConfig, error) {
	authConfig := &AuthConfig{}
	_, err := toml.Decode(content, authConfig)
	if err != nil {
		return nil, err
	}

	return authConfig, nil
}

func filterRegistryScheme(registry string) (string, error) {
	u, err := url.Parse(registry)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(registry, u.Scheme) {
		return strings.TrimPrefix(registry, u.Scheme+"://"), nil
	} else {
		return "", errors.New("registry: address not match after parse")
	}
}
