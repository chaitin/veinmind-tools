package main

import (
	"github.com/chaitin/libveinmind/go/cmd"

	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/authz"
)

var authCmd = &cmd.Command{
	Use:   "authz",
	Short: "Authz as docker plugin",
	RunE: func(cmd *cmd.Command, args []string) error {
		path, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}

		config, err := authz.NewDockerPluginConfig(path)
		if err != nil {
			return err
		}

		options := []authz.DockerServerOption{
			authz.WithPolicy(config.Policies...),
			authz.WithAuthLog(config.Log.AuthZLogPath),
			authz.WithPluginLog(config.Log.PluginLogPath),
			authz.WithListenerUnix(config.Listener.ListenAddr),
		}
		server, err := authz.NewDockerPluginServer(options...)
		if err != nil {
			return err
		}
		runner := authz.NewDefaultRunner(&server)
		return runner.Run()
	},
}

func init() {
	authCmd.Flags().StringP("config", "c", "", "authz config path")
}
