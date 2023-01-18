package registry

import (
	"context"
	"fmt"

	"github.com/gogf/gf/text/gstr"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/pkg/errors"
)

type v2client struct {
	option *Options
}

func NewV2Client(registry string, opts ...Option) (Client, error) {
	option := &Options{
		Insecure: true,
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

	registry, err := name.NewRegistry(c.option.Registry, name.Insecure)
	repos, err := remote.Catalog(ctx, registry, remote.WithAuth(&authn.Basic{
		Username: c.option.Username,
		Password: c.option.Password,
	}))
	if err != nil {
		return nil, errors.Wrap(err, "[registry] can't fetch catalog")
	}

	refs := make([]string, 0)
	for _, repo := range repos {
		ref := repo
		if !gstr.HasPrefix(ref, c.option.Registry) {
			ref = fmt.Sprintf("%s/%s", c.option.Registry, ref)
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
	repositoryOpt = append(repositoryOpt, name.Insecure)
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
