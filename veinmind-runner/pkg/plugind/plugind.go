package plugind

import (
	"context"
	_ "embed"
	"github.com/BurntSushi/toml"
	"log"
)

type PluginServicesConf struct {
	Plugins []PluginConf `toml:"PluginConf"`
}

type PluginConf struct {
	Name     string        `toml:"Name"`
	Services []ServiceConf `toml:"ServiceConf"`
}

type ServiceConf struct {
	Name       string `toml:"Name"`
	ExecScript string `toml:"ExecScript"`
	ExecArgs   string `toml:"ExecArgs"`
	StdoutLog  string `toml:"StdoutLog"`
	StderrLog  string `toml:"StderrLog"`
	Port       string `toml:"Port"`
}

//go:embed conf/service.toml
var conf string

type Plugin struct {
	PluginName string
	Service    []*ServiceRunner
}

type PluginsServices struct {
	Plugins []Plugin
}

var plugind = func() PluginsServices {
	var psConf PluginServicesConf
	_, err := toml.Decode(conf, &psConf)
	if err != nil {
		log.Println(err)
	}
	var ps PluginsServices
	for _, plugin := range psConf.Plugins {
		plugin := initService(plugin)
		if len(plugin.Service) != 0 {
			ps.Plugins = append(ps.Plugins, plugin)
		}
	}
	return ps
}()

func initService(p PluginConf) Plugin {
	var services Plugin
	services.PluginName = p.Name
	for _, service := range p.Services {
		s, err := NewService(service)
		if err != nil {
			log.Println("init error: ", service.Name, err)
		} else {
			services.Service = append(services.Service, s)
		}
	}
	return services
}

func Start() {
	for _, pservice := range plugind.Plugins {
		for _, runner := range pservice.Service {
			if runner != nil {
				runner.Run(context.Background())
				err := runner.Ready()
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func Stop() error {
	for _, service := range plugind.Plugins {
		for _, runner := range service.Service {
			if runner != nil && runner.stop != nil {
				err := runner.Stop()
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func Work() (bool, error) {
	for _, service := range plugind.Plugins {
		for _, runner := range service.Service {
			err := runner.Ready()
			if err != nil {
				return false, err
			}
		}
	}
	return true, nil
}
