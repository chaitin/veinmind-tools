package registry

type Client interface {
	Pull(repo string) (string, error)
	Remove(id string) error
	Auth(config AuthConfig) error
}
