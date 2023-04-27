package main

import (
	"os"
	"strings"

	"github.com/opencontainers/runc/libcontainer"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-unsafe-mount/pkg/engine"
)

var (
	pluginInfo = plugin.Manifest{
		Name:        "veinmind-unsafe-mount",
		Author:      "veinmind-team",
		Description: "detect unsafe mount for container",
	}
	reportService = &report.Service{}
	rootCmd       = &cmd.Command{}
	scanCmd       = &cmd.Command{
		Use: "scan",
	}
	scanContainerCmd = &cmd.Command{
		Use:   "container",
		Short: "scan container command",
	}
)

// scanContainer is func that used to do some action with container
// you can write your plugin scan code here
func scanContainer(c *cmd.Command, container api.Container) error {
	log.Infof("start scan container unsafe mount: %s", container.ID())
	evts, err := engine.DetectContainerUnsafeMount(container)
	if err != nil {
		if !strings.Contains(err.Error(), libcontainer.ErrNotRunning.Error()) { //当扫描所有容器时，会因为有停止的容器而被终止
			log.Error(err.Error())
			return err
		} else {
			log.Infof("skip scanning for containers that are not running for unsafe mount: %s", container.ID())
		}
	}

	for _, evt := range evts {
		err := reportService.Client.Report(&evt)
		if err != nil {
			log.Error(err)
			continue
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.AddCommand(report.MapReportCmd(cmd.MapContainerCommand(scanContainerCmd, scanContainer), reportService))
	rootCmd.AddCommand(cmd.NewInfoCommand(pluginInfo))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
