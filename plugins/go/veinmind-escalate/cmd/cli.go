package main

import (
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-escalate/utils"
	"os"
)

var rootCmd = &cmd.Command{}

var scanCmd = &cmd.Command{
	Use:   "scan",
	Short: "scan mode",
}

var scanImageCmd = &cmd.Command{
	Use:   "image",
	Short: "scan image escalate",
}
var scanContainerCmd = &cmd.Command{
	Use:   "container",
	Short: "scan container escalate",
}

func scanImage(c *cmd.Command, image api.Image) error {
	defer func() {
		err := image.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	err := utils.ImagesScanRun(image)
	if err != nil {
		log.Error(err)
		return err
	}
	// if you want display at runner report, you should send your result to report event
	return nil
}

func scanContainer(c *cmd.Command, container api.Container) error {
	defer func() {
		err := container.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	err := utils.ContainersScanRun(container)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.AddCommand(cmd.MapImageCommand(scanImageCmd, scanImage))
	scanCmd.AddCommand(cmd.MapContainerCommand(scanContainerCmd, scanContainer))

	rootCmd.AddCommand(cmd.NewInfoCommand(plugin.Manifest{
		Name:        "veinmind-escalate",
		Author:      "veinmind-team",
		Description: "detect escalation risk for image&container",
	}))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
