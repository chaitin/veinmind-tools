package conf

import (
	"context"
	"errors"
	"github.com/chaitin/libveinmind/go/plugin/service"
	"golang.org/x/sync/errgroup"
)

const Namespace = "github.com/chaitin/veinmind-tools/veinmind-common/go/service/conf"

type ConfService struct {
	store map[PluginConfNS][]byte
}

type confClient struct {
	ctx   context.Context
	group *errgroup.Group
	Pull  func(ns PluginConfNS) ([]byte, error)
}

func NewConfService() (*ConfService, error) {
	c := new(ConfService)
	c.store = make(map[PluginConfNS][]byte)
	return c, nil
}

func (c *ConfService) Pull(ns PluginConfNS) ([]byte, error) {
	if b, ok := c.store[ns]; ok {
		return b, nil
	} else {
		return nil, errors.New("conf: plugin conf doesn't exist")
	}
}

func (c *ConfService) Store(ns PluginConfNS, b []byte) error {
	c.store[ns] = b
	return nil
}

func (c *ConfService) Add(registry *service.Registry) {
	registry.Define(Namespace, struct{}{})
	registry.AddService(Namespace, "pull", c.Pull)
}
