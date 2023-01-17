package registry

import (
	"context"
)

type Client interface {
	GetRepos(ctx context.Context, opts ...Option) ([]string, error)
	GetRepoTags(ctx context.Context, repo string, opts ...Option) ([]string, error)
}

type Options struct {
	Registry string
	Username string
	Password string
	Token    string
	Insecure bool
}

type Option func(*Options)

func WithLoginByPassword(username, password string) Option {
	return func(o *Options) {
		o.Username = username
		o.Password = password
	}
}

func WithLoginByToken(token string) Option {
	return func(o *Options) {
		o.Token = token
	}
}

func WithInsecure(insecure bool) Option {
	return func(o *Options) {
		o.Insecure = insecure
	}
}
