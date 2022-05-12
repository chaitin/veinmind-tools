package registry

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/distribution/distribution/reference"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/config/configfile"
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

type RegistryDockerClient struct {
	ctx     context.Context
	auth    map[string]Auth
	options []remote.Option
}

func parseDockerAuthConfig(path string) (map[string]Auth, error) {
	dockerConfig := configfile.ConfigFile{}
	authConfigMap := make(map[string]Auth)

	if _, err := os.Stat(path); err != nil {
		dockerConfigByte, err := ioutil.ReadFile(path)

		err = json.Unmarshal(dockerConfigByte, &dockerConfig)
		if err != nil {
			return nil, err
		} else {
			for server, config := range dockerConfig.AuthConfigs {
				u, err := url.Parse(server)
				registryName := ""
				if err != nil {
					registryName = server
				} else {
					registryName = u.Host
				}

				registry, err := name.NewRegistry(registryName)
				if err != nil {
					log.Error(err)
					continue
				}

				if config.Auth != "" {
					authDecode, err := base64.StdEncoding.DecodeString(config.Auth)
					if err == nil {
						authSplit := strings.Split(string(authDecode), ":")
						if len(authSplit) == 2 {
							auth := Auth{
								Username: authSplit[0],
								Password: authSplit[1],
							}
							authConfigMap[registry.String()] = auth
						} else {
							log.Error("docker config auth block length wrong")
							continue
						}
					} else {
						log.Error(err)
						continue
					}
				}
			}
			return authConfigMap, nil
		}
	} else {
		return nil, err
	}
}

func NewRegistryDockerClient(opts ...Option) (Client, error) {
	c := &RegistryDockerClient{}
	c.ctx = context.Background()
	c.auth = make(map[string]Auth)

	// Get Auth Token From Config File
	auth, err := parseDockerAuthConfig(dockerConfigPath)
	if err != nil {
		log.Error(err)
	} else {
		c.auth = auth
	}

	// Options handle
	for _, opt := range opts {
		cNew, err := opt(c)
		if err != nil {
			log.Error(err)
			continue
		}
		c = cNew.(*RegistryDockerClient)
	}

	var clientOpts []remote.Option
	clientOpts = append(clientOpts, remote.WithTransport(&http.Transport{
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
	c.options = clientOpts

	return c, nil
}

func (client *RegistryDockerClient) GetRepo(repo string, options ...remote.Option) (*remote.Descriptor, error) {
	options = append(options, client.options...)
	named, err := reference.ParseDockerRef(repo)
	if err != nil {
		return nil, err
	}

	domain := reference.Domain(named)
	var auth Auth
	if v, ok := client.auth[domain]; ok {
		auth = v
	}

	if auth.Username != "" && auth.Password != "" {
		options = append(options, remote.WithAuth(&authn.Basic{
			Username: auth.Username,
			Password: auth.Password,
		}))
	}

	ref, err := name.ParseReference(repo)
	if err != nil {
		return nil, err
	}
	return remote.Get(ref, options...)
}

func (client *RegistryDockerClient) GetRepoTags(repo string, options ...remote.Option) ([]string, error) {
	options = append(options, client.options...)
	named, err := reference.ParseDockerRef(repo)
	if err != nil {
		return nil, err
	}

	domain := reference.Domain(named)
	var auth Auth
	if v, ok := client.auth[domain]; ok {
		auth = v
	}

	if auth.Username != "" && auth.Password != "" {
		options = append(options, remote.WithAuth(&authn.Basic{
			Username: auth.Username,
			Password: auth.Password,
		}))
	}

	repoR, err := name.NewRepository(repo)
	if err != nil {
		return nil, err
	}
	return remote.List(repoR, options...)
}

func (client *RegistryDockerClient) GetRepos(address string, options ...remote.Option) (repos []string, err error) {
	options = append(options, client.options...)
	var auth Auth
	if v, ok := client.auth[address]; ok {
		auth = v
	}

	if auth.Username != "" && auth.Password != "" {
		options = append(options, remote.WithAuth(&authn.Basic{
			Username: auth.Username,
			Password: auth.Password,
		}))
	}

	regsitry, err := name.NewRegistry(address)
	if err != nil {
		return nil, err
	}

	last := ""
	for {
		reposTemp := []string{}
		reposTemp, err = remote.CatalogPage(regsitry, last, 10000, options...)
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

func (client *RegistryDockerClient) Auth(config AuthConfig) error {
	for _, auth := range config.Auths {
		client.auth[auth.Registry] = auth
	}

	return nil
}

func (client *RegistryDockerClient) Pull(repo string) (string, error) {
	c, err := dockercli.NewClientWithOpts(dockercli.FromEnv, dockercli.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	named, err := reference.ParseDockerRef(repo)
	if err != nil {
		return "", err
	}

	domain := reference.Domain(named)
	var auth Auth
	if v, ok := client.auth[domain]; ok {
		auth = v
	}

	// Generate Auth Token
	token, err := command.EncodeAuthToBase64(dockertypes.AuthConfig{
		Username: auth.Username,
		Password: auth.Password})

	var closer io.ReadCloser
	if token == "" {
		closer, err = c.ImagePull(client.ctx, repo, dockertypes.ImagePullOptions{})
	} else {
		closer, err = c.ImagePull(client.ctx, repo, dockertypes.ImagePullOptions{
			RegistryAuth: token,
		})
	}

	_, err = ioutil.ReadAll(closer)
	if err != nil {
		return "", err
	}

	return named.String(), nil
}

func (client *RegistryDockerClient) Remove(id string) error {
	c, err := dockercli.NewClientWithOpts(dockercli.FromEnv, dockercli.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	_, err = c.ImageRemove(client.ctx, id, dockertypes.ImageRemoveOptions{
		Force:         true,
		PruneChildren: true,
	})
	return err
}
