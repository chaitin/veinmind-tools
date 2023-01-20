package plugind

import (
	"context"
	_ "embed"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/log"
	"golang.org/x/sync/errgroup"
)

//go:embed conf/service.toml
var config string

func NewManager() (*Manager, error) {
	var pluginsManager Manager
	_, err := toml.Decode(config, &pluginsManager)
	if err != nil {
		return nil, err
	}

	return &pluginsManager, nil
}

func (c *Manager) StartWithContext(ctx context.Context, name string) error {
	for _, plugin := range c.Plugins {
		if plugin.Name != name {
			continue
		}
		for _, s := range plugin.Service {
			if !s.Running {
				log.GetModule(log.PlugindModuleKey).Infof("%s related service %s try start", plugin.Name, s.Name)
				err := svcManager.Start(ctx, s)
				if err != nil {
					return err
				}
				log.GetModule(log.PlugindModuleKey).Infof("%s related service %s success start", plugin.Name, s.Name)
			}
		}
	}
	return nil
}

func (c *Manager) Wait() {
	svcManager.wg.Wait()
}

type serviceManager struct {
	wg *sync.WaitGroup
}

func newServiceManager() serviceManager {
	return serviceManager{
		wg: &sync.WaitGroup{},
	}
}

var svcManager = newServiceManager()

func (s *serviceManager) Start(ctx context.Context, conf *Service) error {
	checkChains := make([]serviceCheckFunc, 0)
	for _, check := range conf.Check {
		fn, ok := serviceChecks[check.Type]
		if ok {
			checkChains = append(checkChains, fn(check.Value))
		}
	}

	options := make([]serviceOption, 0)
	options = append(options, withStdout(conf.Stdout), withStderr(conf.Stderr))
	options = append(options, withTimeout(conf.Timeout), withWaitGroup(s.wg))
	options = append(options, withCheckChains(checkChains...))

	svc := newService(ctx, conf.Command, options...)
	if err := svc.start(); err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(svc.ready)
	if err := g.Wait(); err != nil {
		return err
	}

	conf.Running = true

	go svc.daemon()

	return nil
}
