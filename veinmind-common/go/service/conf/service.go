package conf

import (
	"context"
	"errors"
	"github.com/chaitin/libveinmind/go/plugin/service"
	"golang.org/x/sync/errgroup"
)

const Namespace = "github.com/chaitin/veinmind-tools/veinmind-common/go/service/conf"

type ConfService struct {
	store map[string][]byte
}

type confClient struct {
	ctx    context.Context
	group  *errgroup.Group
	Pull func(ns PluginConfNS) ([]byte, error)
}

func (c *ConfService) Pull(pluginName string) ([]byte, error){
	if b, ok := c.store[pluginName]; ok {
		return b, nil
	}else{
		return nil, errors.New("conf: plugin conf doesn't exist")
	}
}

func (c *ConfService) Store(pluginName string, b []byte) error {
	c.store[pluginName] = b
	return nil
}

func (c *ConfService) Add(registry *service.Registry)  {
	registry.Define(Namespace, struct {}{})
	registry.AddService(Namespace, "pull", c.Pull)
}
