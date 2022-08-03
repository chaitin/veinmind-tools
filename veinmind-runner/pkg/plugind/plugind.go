package plugind

import (
	"context"
	"errors"
	"github.com/BurntSushi/toml"
	"sync"
)

func NewPlugindConfig(config string) (*Config, error) {
	var plugindConfig Config
	_, err := toml.DecodeFile(config, &plugindConfig)
	if err != nil {
		return nil, err
	}
	return &plugindConfig, nil
}

func (c *Config) StartWithContext(ctx context.Context, name string) error {
	for _, conf := range c.Plugin {
		if conf.Name == name {
			for _, s := range conf.Service {
				err := pluginSvc.Start(ctx, s)
				if err != nil {
					return err
				}
			}
			return nil
		}
	}
	return errors.New("")
}

func (c *Config) Wait() {
	pluginSvc.wg.Wait()
}

type pluginService struct {
	wg       *sync.WaitGroup
	services sync.Map
}

func newPluginService() pluginService {
	return pluginService{
		wg:       &sync.WaitGroup{},
		services: sync.Map{},
	}
}

var pluginSvc = newPluginService()

func (s *pluginService) Start(ctx context.Context, conf *ServiceConf) error {
	checkChains := make([]serviceCheckFunc, 0)
	for _, check := range conf.Check {
		fn, ok := serviceChecks[check.Type]
		if ok {
			checkChains = append(checkChains, fn(check.Value))
		}
	}

	options := make([]serviceOption, 0)
	options = append(options, withStdout(conf.Stdout), withStderr(conf.Stderr))
	options = append(options, withTimeout(conf.Timeout))
	options = append(options, withCheckChains(checkChains...))

	svc := newService(ctx, conf.Command, options...)
	if err := svc.Start(s.wg); err != nil {
		return err
	}
	s.services.Store(conf.Name, svc)

	return nil
}
