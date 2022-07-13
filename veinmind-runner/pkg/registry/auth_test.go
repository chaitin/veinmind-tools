package registry

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseAuthConfig(t *testing.T) {
	content := `[[auths]]
	registry = "index.docker.io"
	username = "admin"
	password = "password"
	[[auths]]
	registry = "private.net"
	username = "admin"
	password = "password"
	`
	config, err := parseAuthConfigFromString(content)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, &AuthConfig{Auths: []Auth{
		{
			Registry: "index.docker.io",
			Username: "admin",
			Password: "password",
		},
		{
			Registry: "private.net",
			Username: "admin",
			Password: "password",
		},
	}}, config)
}
