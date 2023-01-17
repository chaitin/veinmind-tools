package main

import (
	"context"
	"encoding/json"

	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
)

// listCmd used to display relevant information
var listCmd = &cmd.Command{
	Use:   "list",
	Short: "list relevant information",
}

var listPluginCmd = &cmd.Command{
	Use:   "plugin",
	Short: "list plugin information",
	RunE: func(cmd *cmd.Command, args []string) error {
		ps, err := plugin.DiscoverPlugins(context.Background(), ".")
		if err != nil {
			return err
		}

		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			return err
		}

		for _, p := range ps {
			if verbose {
				pJsonByte, err := json.MarshalIndent(p.Manifest, "", "	")
				if err != nil {
					log.Error(err)
					continue
				}
				log.Info("\n" + string(pJsonByte))
			} else {
				log.Info("Plugin Name: " + p.Name)
			}
		}
		return nil
	},
}

// perhaps let user display load IaC rules
var listIaCRuleCmd = &cmd.Command{}

func init() {
	listCmd.AddCommand(listPluginCmd)

	listPluginCmd.Flags().BoolP("verbose", "v", false, "verbose mode")
}
