package main

import (
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/prometheus/common/log"
)

var scanCmd = &cmd.Command{
	Use:   "scan",
	Short: "scan command",
}

func scan(_ *cmd.Command, image api.Image) error {
	log.Info(image.ID())
	return nil
}

func main() {
	cmd.MapImageCommand(scanCmd, scan)
	cmd.NewInfoCommand(plugin.Manifest{
		Name:   "veinmind-example",
		Author: "veinmind-team",
	})
	if err := scanCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
