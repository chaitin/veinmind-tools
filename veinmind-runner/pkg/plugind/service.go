package plugind

import (
	"errors"
	"github.com/BurntSushi/toml"
	"log"
	"os/exec"
)

type PluginService struct {
	Plugins []Plugin `toml:"Plugin"`
}

type Plugin struct {
	Name     string    `toml:"Name"`
	Services []Service `toml:"Service"`
}

type Service struct {
	Name       string `toml:"Name"`
	ExecScript string `toml:"ExecScript"`
	ExecNeed   bool   `toml:"ExecNeed"`
	ExecStop   bool   `toml:"ExecStop"`
}

var servicePath = "./conf/service.toml"

var plugind = func() PluginService {
	var pluginServices PluginService
	_, err := toml.DecodeFile(servicePath, &pluginServices)
	if err != nil {
		log.Println(err)
	}
	return pluginServices
}()

func Run(s Service, action string) error {
	_, err := exec.Command(s.ExecScript, action).Output() //nolint:gosec
	if err != nil {
		return err
	}
	return nil
}

func StartService(s Service) error {
	if s.ExecNeed {
		return Run(s, "start")
	}
	return nil
}

func StopService(s Service) error {
	if s.ExecStop {
		return Run(s, "stop")
	}
	return nil
}

func StatusService(s Service) bool {
	if err := Run(s, "status"); err != nil {
		return false
	}
	return true
}

func StartServices(s []Service) error {
	for _, service := range s {
		if !StatusService(service) {
			err := StartService(service)
			if err != nil {
				return errors.New("Start " + service.Name + " error: " + err.Error())
			}
		}
	}
	return nil
}

func StopServices(s []Service) error {
	for _, service := range s {
		if StatusService(service) {
			err := StopService(service)
			if err != nil {
				return errors.New("Stop " + service.Name + " error: " + err.Error())
			}
		}
	}
	return nil
}

func StatusServices(s []Service) bool {
	var isAllWorking = true
	for _, service := range s {
		isAllWorking = isAllWorking && StatusService(service)
	}
	return isAllWorking
}

func StartPluginsService() error {
	for _, plugin := range plugind.Plugins {
		return StartServices(plugin.Services)
	}
	return nil
}

func StopPluginsService() error {
	for _, plugin := range plugind.Plugins {
		return StopServices(plugin.Services)
	}
	return nil
}

func StatusPluginsServices(pluginName string) bool {
	for _, plugin := range plugind.Plugins {
		if plugin.Name == pluginName {
			return StatusServices(plugin.Services)
		}
	}
	return false
}

func StatusAllPluginsService() bool {
	workingFlag := true
	for _, plugin := range plugind.Plugins {
		workingFlag = workingFlag && StatusServices(plugin.Services)
	}
	return workingFlag
}
