package registry

import (
	"context"
	"fmt"

	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/gogf/gf/text/gstr"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

type v2client struct {
	option *Options
}

func NewV2Client(registry string, opts ...Option) (Client, error) {
	option := &Options{}

	// trans: http://xxxx || https://xxxx => xxxx
	// && change options insecure
	if gstr.HasPrefix("http://", registry) {
		registry = gstr.Replace(registry, "http://", "", -1)
		option.Insecure = false
	}

	if gstr.HasPrefix("https://", registry) {
		registry = gstr.Replace(registry, "https://", "", -1)
		option.Insecure = true
	}

	// user's options > scheme
	for _, o := range opts {
		o(option)
	}

	option.Registry = registry
	c := &v2client{
		option: option,
	}

	return c, nil
}

func (c *v2client) GetRepos(ctx context.Context, opts ...Option) ([]string, error) {

	for _, o := range opts {
		o(c.option)
	}

	repositoryOpt := make([]name.Option, 0)
	if c.option.Insecure {
		repositoryOpt = append(repositoryOpt, name.Insecure)
	}
	// try strict first
	if reference, err := name.ParseReference(c.option.Registry, append(repositoryOpt, name.StrictValidation)...); err == nil {
		return []string{reference.Name()}, nil
	}

	reference, err := name.ParseReference(c.option.Registry, repositoryOpt...)
	if err != nil {
		return nil, err
	}

	// dockerhub: just call remote
	if reference.Context().RegistryStr() == "index.docker.io" {
		log.Warnf("found server: docker")
		log.Warnf("Currently, docker.io authentication is not supported, so it is automatically scanned as a public image without authentication information")
		return []string{reference.Name()}, nil
	}

	// others: list all image reference
	repos, err := remote.Catalog(ctx, reference.Context().Registry, remote.WithAuth(&authn.Basic{
		Username: c.option.Username,
		Password: c.option.Password,
	}))

	if err != nil {
		return []string{}, err
	}

	refs := make([]string, 0)
	for _, repo := range repos {
		ref := repo
		if !gstr.HasPrefix(ref, reference.Context().RegistryStr()) {
			ref = fmt.Sprintf("%s/%s", reference.Context().RegistryStr(), ref)
		}
		refs = append(refs, ref)
	}

	return refs, nil
}

func (c *v2client) GetRepoTags(ctx context.Context, repo string, opts ...Option) ([]string, error) {
	for _, o := range opts {
		o(c.option)
	}
	repositoryOpt := make([]name.Option, 0)
	if c.option.Insecure {
		repositoryOpt = append(repositoryOpt, name.Insecure)
	}
	reference, err := name.ParseReference(repo, repositoryOpt...)
	if err != nil {
		return []string{}, err
	}
	repos := reference.Context()

	return remote.List(repos, remote.WithAuth(&authn.Basic{
		Username: c.option.Username,
		Password: c.option.Password,
	}))
}
