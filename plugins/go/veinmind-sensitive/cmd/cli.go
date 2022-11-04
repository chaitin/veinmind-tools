package main

import (
	"context"

	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-sensitive/rule"
	"github.com/spf13/cobra"
)

var (
	rootCommand = &cmd.Command{}
	scanCommand = &cmd.Command{
		Use:   "scan image sensitive info",
		Short: "scan image sensitive info",
		PreRun: func(cmd *cobra.Command, args []string) {
			rule.Init()
		},
	}
)

func init() {
	rootCommand.AddCommand(cmd.MapImageCommand(scanCommand, Scan))
	rootCommand.AddCommand(cmd.NewInfoCommand(plugin.Manifest{
		Name:        "veinmind-sensitive",
		Author:      "veinmind-team",
		Description: "veinmind-sensitive scan image sensitive data",
	}))
}

func main() {
	if err := rootCommand.ExecuteContext(context.Background()); err != nil {
		panic(err)
	}
}
