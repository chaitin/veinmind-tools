package plugind

import (
	"context"
	_ "embed"
	"errors"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/plugind/service"
	"reflect"
)

func RunPluginService(ctx context.Context, Name string) error {
	for _, p := range plugind {
		if p.PluginName == Name {
			ctx, p.StopDaemon = context.WithCancel(ctx)
			err := p.runService(ctx)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("can not find the plugins service")
}

// runService start the Plugin's all Services runner
func (p *Plugin) runService(ctx context.Context) error {
	p.daemon(p.creatCase(ctx))
	err := p.startAllService()
	if err != nil {
		return err
	}
	return nil
}

func (p *Plugin) creatCase(ctx context.Context) []reflect.SelectCase {
	var cases []reflect.SelectCase
	for _, runner := range p.Service {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(runner.Signal),
		})
	}
	cases = append(cases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(ctx.Done()),
	})
	return cases
}

//daemon will create a coroutine user to monitor
//whether the already running service exits.
//If it exits, it will be pulled up again.
func (p *Plugin) daemon(cases []reflect.SelectCase) {
	go func() {
		p.syncFlag.Add(1)
		defer p.syncFlag.Done()
		for {
			_, rev, ok := reflect.Select(cases)
			if !rev.IsValid() || !ok {
				p.killService()
				return
			}
			err := p.restartService(rev.String())
			if err != nil {
				log.Error(err)
			}
		}
	}()
}

// startAllService start all the runner
func (p *Plugin) startAllService() error {
	for _, runner := range p.Service {
		err := runner.Start()
		if err != nil {
			p.StopDaemon()
			return err
		}
		p.RunnerMap.Store(runner.Uuid, runner)
	}
	return nil
}

// restartService the process according to the Uuid of the process
func (p *Plugin) restartService(uuid string) error {
	value, ok := p.RunnerMap.Load(uuid)
	if ok {
		return value.(*service.Runner).Start()
	}
	return nil
}

// killService all surviving processes
func (p *Plugin) killService() {
	p.RunnerMap.Range(func(key, value any) bool {
		p.RunnerMap.Delete(key)
		err := value.(*service.Runner).Stop()
		if err != nil {
			log.Error(err)
		}
		return true
	})
}
