package plugind

import (
	_ "embed"
	"github.com/BurntSushi/toml"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/plugind/service"
	"sync"
)

//go:embed conf/service.toml
var servicePath string

type Conf struct {
	Plugins []PluginConf `toml:"PluginConf"`
}

type PluginConf struct {
	Name     string         `toml:"Name"`
	Services []service.Conf `toml:"ServiceConf"`
}

type Plugin struct {
	PluginName string
	Service    []*service.Runner
	RunnerMap  *sync.Map
	StopDaemon func()
	syncFlag   *sync.WaitGroup
}

var plugind = func() []Plugin {
	var psConf Conf
	var ps []Plugin

	_, err := toml.Decode(servicePath, &psConf)
	if err != nil {
		return ps
	}

	for _, plugin := range psConf.Plugins {
		plugin := initService(plugin)
		if len(plugin.Service) != 0 {
			ps = append(ps, plugin)
		}
	}
	return ps
}()

func initService(plugin PluginConf) Plugin {
	p := Plugin{
		PluginName: plugin.Name,
		RunnerMap:  &sync.Map{},
		syncFlag:   &sync.WaitGroup{},
	}
	for _, s := range plugin.Services {
		runner, err := service.NewRunner(s)
		if err != nil {
			continue
		}
		p.Service = append(p.Service, runner)
	}
	return p
}
