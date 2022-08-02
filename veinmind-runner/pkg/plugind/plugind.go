package plugind

import (
	"context"
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

func StartWithContext(ctx context.Context, conf []*ServiceConf) error {
	for _, c := range conf {
		err := pluginSvc.Start(ctx, c)
		if err != nil {
			return err
		}
	}
	return nil
}

type pluginService struct {
	services sync.Map
}

func newPluginService() pluginService {
	return pluginService{
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
	if err := svc.Start(); err != nil {
		return err
	}
	s.services.Store(conf.Name, svc)

	return nil
}

func (s *pluginService) IsAlive(name string) (bool, error) {
	val, ok := s.services.Load(name)
	if !ok {
		return false, SvcNotExist
	}
	svc := val.(service)

	return svc.IsAlive(), nil
}
