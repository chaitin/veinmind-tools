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

var scanImageCmd = &cmd.Command{
	Use:   "image",
	Short: "scan image escalate",
}
var scanContainerCmd = &cmd.Command{
	Use:   "container",
	Short: "scan container escalate",
}

// scanImage is func that used to do some action with Images
// you can write your plugin scan code here
func scanImage(c *cmd.Command, image api.Image) error {
	defer func() {
		err := image.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	// do something here
	log.Info(image.ID())
	utils.ImagesScanRun(image)

	// if you want display at runner report, you should send your result to report event
	return nil
}

// scanContainer is func that used to do some action with container
// you can write your plugin scan code herex
func scanContainer(c *cmd.Command, container api.Container) error {

	// do something here
	log.Info(container.ID())
	utils.ContainersScanRun(container)
	// if you want display at runner report, you should send your result to report event
	return nil
}

func init() {
	rootCmd.AddCommand(cmd.MapImageCommand(scanImageCmd, scanImage))
	rootCmd.AddCommand(cmd.MapContainerCommand(scanContainerCmd, scanContainer))
	rootCmd.AddCommand(cmd.NewInfoCommand(plugin.Manifest{
		Name:        "veinmind-escalate",
		Author:      "veinmind-team",
		Description: "veinmind-escalate Description",
	}))
	//scanCmd.AddCommand(cmd.MapImageCommand(scanImageCmd, scanImage))
	//scanCmd.AddCommand(cmd.MapContainerCommand(scanContainerCmd, scanContainer))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
