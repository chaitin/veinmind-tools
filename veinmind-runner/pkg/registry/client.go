package registry

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/config/types"
	dockertypes "github.com/docker/docker/api/types"
	dockercli "github.com/docker/docker/client"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const dockerConfigPath = "/root/.docker/config.json"

type Auth struct {
	Username string
	Password string
}

type RegistryClient struct {
	ctx     context.Context
	target  name.Registry
	auth    *Auth
	options []remote.Option
}

func NewRegistryClient(registryAddress string, auth *Auth, opts ...remote.Option) (*RegistryClient, error) {
	c := &RegistryClient{}
	c.ctx = context.Background()
	target, err := name.NewRegistry(registryAddress)
	if err != nil {
		return nil, err
	}
	c.target = target
	if auth == nil {
		auth = &Auth{}

		// Get Auth Token From Config File
		dockerConfig := configfile.ConfigFile{}
		if _, err := os.Stat(dockerConfigPath); !os.IsNotExist(err) {
			dockerConfigByte, err := ioutil.ReadFile(dockerConfigPath)
			if err == nil {
				err = json.Unmarshal(dockerConfigByte, &dockerConfig)
				if err == nil {
					authConfig := types.AuthConfig{}
					if config, ok := dockerConfig.AuthConfigs[c.target.Name()]; ok {
						authConfig = config
					} else {
						for server, config := range dockerConfig.AuthConfigs {
							if !strings.Contains(server, "://") {
								server = "//" + server
							}
							u1, err1 := url.Parse(server)
							u2, err2 := url.Parse("//" + c.target.Name())
							if err1 != nil {
								log.Error(err)
								continue
							}
							if err2 != nil {
								log.Error(err)
								continue
							}
							if strings.EqualFold(u1.Host, u2.Host) {
								authConfig = config
							}
						}
					}
					if authConfig.Auth != "" {
						authDecode, err := base64.StdEncoding.DecodeString(authConfig.Auth)
						if err == nil {
							authSplit := strings.Split(string(authDecode), ":")
							if len(authSplit) == 2 {
								auth.Username = authSplit[0]
								auth.Password = authSplit[1]
							} else {
								log.Error("docker config auth block length wrong")
							}
						} else {
							log.Error(err)
						}
					} else if authConfig.Username != "" && authConfig.Password != "" {
						auth.Username = authConfig.Username
						auth.Password = authConfig.Password
					}
				} else {
					log.Error(err)
				}
			} else {
				log.Error(err)
			}
		}
		c.auth = auth
	} else {
		c.auth = auth
	}

	opts = append(opts, remote.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}))
	c.options = opts

	return c, nil
}

func (client *RegistryClient) GetRepo(repo string, options ...remote.Option) (*remote.Descriptor, error) {
	options = append(options, client.options...)
	if client.auth.Username != "" && client.auth.Password != "" {
		options = append(options, remote.WithAuth(&authn.Basic{
			Username: client.auth.Username,
			Password: client.auth.Password,
		}))
	}

	ref, err := name.ParseReference(repo)
	if err != nil {
		return nil, err
	}
	return remote.Get(ref, options...)
}

func (client *RegistryClient) GetRepoTags(repo string, options ...remote.Option) ([]string, error) {
	options = append(options, client.options...)
	if client.auth.Username != "" && client.auth.Password != "" {
		options = append(options, remote.WithAuth(&authn.Basic{
			Username: client.auth.Username,
			Password: client.auth.Password,
		}))
	}

	repoR, err := name.NewRepository(repo)
	if err != nil {
		return nil, err
	}
	return remote.List(repoR, options...)
}

func (client *RegistryClient) GetRepos(options ...remote.Option) (repos []string, err error) {
	options = append(options, client.options...)
	if client.auth.Username != "" && client.auth.Password != "" {
		options = append(options, remote.WithAuth(&authn.Basic{
			Username: client.auth.Username,
			Password: client.auth.Password,
		}))
	}

	last := ""

	for {
		reposTemp := []string{}
		reposTemp, err = remote.CatalogPage(client.target, last, 10000, options...)
		if err != nil {
			break
		}

		if len(reposTemp) > 0 {
			repos = append(repos, reposTemp...)
		} else {
			break
		}

		last = reposTemp[len(reposTemp)-1]
	}

	return repos, err
}

func (client *RegistryClient) Pull(repo string) (io.ReadCloser, error) {
	c, err := dockercli.NewClientWithOpts(dockercli.FromEnv, dockercli.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	// Generate Auth Token
	token := ""
	if client.auth != nil {
		token, err = command.EncodeAuthToBase64(dockertypes.AuthConfig{
			Username: client.auth.Username,
			Password: client.auth.Password,
		})
	}

	if token == "" {
		return c.ImagePull(client.ctx, repo, dockertypes.ImagePullOptions{})
	} else {
		return c.ImagePull(client.ctx, repo, dockertypes.ImagePullOptions{
			RegistryAuth: token,
		})
	}
}

func (client *RegistryClient) Remove(id string) ([]dockertypes.ImageDeleteResponseItem, error) {
	c, err := dockercli.NewClientWithOpts(dockercli.FromEnv, dockercli.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return c.ImageRemove(client.ctx, id, dockertypes.ImageRemoveOptions{
		Force:         true,
		PruneChildren: true,
	})
}
