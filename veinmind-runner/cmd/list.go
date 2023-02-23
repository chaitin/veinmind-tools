package main

import (
	"context"
	"encoding/json"

	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"

	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/log"
)

// listCmd used to display relevant information
var listCmd = &cmd.Command{
	Use:   "list",
	Short: "List relevant information",
}

var listPluginCmd = &cmd.Command{
	Use:   "plugin",
	Short: "List plugin information",
	RunE: func(cmd *cmd.Command, args []string) error {
		ps, err := plugin.DiscoverPlugins(context.Background(), "./plugin")
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
					log.GetModule(log.CmdModuleKey).Error(err)
					continue
				}
				log.GetModule(log.CmdModuleKey).Info("\n" + string(pJsonByte))
			} else {
				log.GetModule(log.CmdModuleKey).Infof("plugin name: %s", p.Name)
			}
		}
		return nil
	},
}

func init() {
	listCmd.AddCommand(listPluginCmd)

	listPluginCmd.Flags().BoolP("verbose", "v", false, "verbose mode")
}
