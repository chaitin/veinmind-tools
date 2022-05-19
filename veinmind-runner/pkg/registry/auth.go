package registry

import "github.com/BurntSushi/toml"

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
