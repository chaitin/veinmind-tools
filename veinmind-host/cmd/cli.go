package main

import (
	"context"
	"github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/libveinmind/go/plugin/service"
	"os"
	"path"
)

var rootCmd = &cmd.Command{}
var scanCmd = &cmd.Command{
	Use:   "scan",
	Short: "perform hosted scan command",
}

func scan(c *cmd.Command, image api.Image) error {
	// Discover Plugins
	ctx := c.Context()
	log.Info("Start discovering")
	glob, err := c.Flags().GetString("glob")
	var ps []*plugin.Plugin
	if err == nil && glob != "" {
		ps, err = plugin.DiscoverPlugins(ctx, ".", plugin.WithGlob(glob))
	}else{
		ps, err = plugin.DiscoverPlugins(ctx, ".")
	}
	if err != nil {
		return err
	}
	for _, p := range ps {
		log.Infof("Discovered plugin: %#v\n", p.Name)
	}
	log.Info("Done discovering")


	log.Infof("Start scanning")
	if err := cmd.ScanImages(ctx, ps, []api.Image{image},
		plugin.WithExecInterceptor(func(
			ctx context.Context, plug *plugin.Plugin, c *plugin.Command,
			next func(context.Context, ...plugin.ExecOption) error,
		) error {
			reg := service.NewRegistry()
			reg.AddServices(log.WithFields(log.Fields{
				"plugin":  plug.Name,
				"command": path.Join(c.Path...),
			}))
			return next(ctx, reg.Bind())
		})); err != nil {
		return err
	}
	log.Infof("Done scanning")
	return nil
}

func init() {
	rootCmd.AddCommand(	cmd.MapImageCommand(scanCmd, scan))
	scanCmd.Flags().StringP("glob", "g", "", "specifies the pattern of plugin file to find")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
