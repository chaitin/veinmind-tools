package main

import (
	"context"

	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/veinmind-common-go/service/report"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-sensitive/rule"
)

var (
	PluginInfo = plugin.Manifest{
		Name:        "veinmind-sensitive",
		Author:      "veinmind-team",
		Description: "veinmind-sensitive scan image sensitive data",
		Version:     "v1.1.5",
	}

	reportService = &report.Service{}
	rootCommand   = &cmd.Command{}
	scanCmd       = &cmd.Command{
		Use: "scan",
	}
	scanImageCmd = &cmd.Command{
		Use:   "image",
		Short: "scan image sensitive info",
	}
)

func init() {
	rule.Init()
	rootCommand.AddCommand(scanCmd)
	scanCmd.AddCommand(report.MapReportCmd(cmd.MapImageCommand(scanImageCmd, Scan), reportService))
	rootCommand.AddCommand(cmd.NewInfoCommand(PluginInfo))
}

func main() {
	if err := rootCommand.ExecuteContext(context.Background()); err != nil {
		panic(err)
	}
}
