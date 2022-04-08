package main

import (
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"os"
)

var rootCmd = &cmd.Command{}
var scanCmd = &cmd.Command{
	Use:   "scan",
	Short: "scan command",
}

func scan(c *cmd.Command, image api.Image) error {
	log.Info(image.ID())
	return nil
}

func init() {
	rootCmd.AddCommand(cmd.MapImageCommand(scanCmd, scan))
	rootCmd.AddCommand(cmd.NewInfoCommand(plugin.Manifest{
		Name:   "veinmind-example",
		Author: "veinmind-team",
	}))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
