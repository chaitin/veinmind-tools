package plugind

import (
	_ "embed"
	"github.com/BurntSushi/toml"
	"os"
	"os/exec"
	"sync"
	"time"
)

//go:embed conf/service.toml
var servicePath string

var (
	Signal    chan string
	RunnerMap sync.Map
)

type Conf struct {
	Plugins []PluginConf `toml:"PluginConf"`
}

type ServiceConf struct {
	Name      string `toml:"Name"`
	Command   string `toml:"Command"`
	StdoutLog string `toml:"StdoutLog"`
	StderrLog string `toml:"StderrLog"`
	Port      string `toml:"Port"`
}

type PluginConf struct {
	Name     string        `toml:"Name"`
	Services []ServiceConf `toml:"ServiceConf"`
}

type Plugin struct {
	PluginName string
	Service    []*Runner
}

type Runner struct {
	Name    string
	Uuid    string
	Command string
	Stderr  *os.File
	Stdout  *os.File
	Port    string
	Cmd     *exec.Cmd
	TimeOut time.Duration
}

func NewPluginServices() ([]Plugin, error) {
	var psConf Conf
	var ps []Plugin

	Signal = make(chan string)
	_, err := toml.Decode(servicePath, &psConf)
	if err != nil {
		return ps, err
	}

	for _, plugin := range psConf.Plugins {
		plugin := initService(plugin)
		if len(plugin.Service) != 0 {
			ps = append(ps, plugin)
		}
	}

	return ps, nil
}

func initService(plugin PluginConf) Plugin {
	p := Plugin{
		PluginName: plugin.Name,
	}
	for _, service := range plugin.Services {
		runner, err := newRunner(service)
		if err != nil {
			continue
		}
		p.Service = append(p.Service, runner)
	}
	return p
}
