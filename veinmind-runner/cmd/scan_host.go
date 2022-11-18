package main

import (
	"context"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/libveinmind/go/plugin/service"
	"path"
)

var scanHostCmd = &cmd.Command{
	Use:   "scan-host",
	Short: "perform hosted scan command",
}

// scan-host command
var scanHostImageCmd = &cmd.Command{
	Use:      "image",
	Short:    "perform hosted image scan",
	PreRunE:  scanPreRunE,
	PostRunE: scanPostRunE,
}
var scanHostContainerCmd = &cmd.Command{
	Use:      "container",
	Short:    "perform hosted container scan",
	PreRunE:  scanPreRunE,
	PostRunE: scanPostRunE,
}

func scanImage(c *cmd.Command, image api.Image) error {
	refs, err := image.RepoRefs()
	ref := ""
	if err == nil && len(refs) > 0 {
		ref = refs[0]
	} else {
		ref = image.ID()
	}

	// Get threads value
	t, err := c.Flags().GetInt("threads")
	if err != nil {
		t = 5
	}

	log.Infof("Scan image: %#v\n", ref)
	if err := cmd.ScanImage(ctx, ps, image,
		plugin.WithExecInterceptor(func(
			ctx context.Context, plug *plugin.Plugin, c *plugin.Command, next func(context.Context, ...plugin.ExecOption) error,
		) error {
			// Register Service
			reg := service.NewRegistry()
			reg.AddServices(log.WithFields(log.Fields{
				"plugin":  plug.Name,
				"command": path.Join(c.Path...),
			}))
			reg.AddServices(reportService)

			// Next Plugin
			return next(ctx, reg.Bind())
		}), plugin.WithExecParallelism(t)); err != nil {
		return err
	}
	return nil
}

func scanContainer(c *cmd.Command, container api.Container) error {

	ref := container.Name()

	// Get threads value
	t, err := c.Flags().GetInt("threads")
	if err != nil {
		t = 5
	}

	log.Infof("Scan container: %#v\n", ref)
	if err := cmd.ScanContainer(ctx, ps, container,
		plugin.WithExecInterceptor(func(
			ctx context.Context, plug *plugin.Plugin, c *plugin.Command, next func(context.Context, ...plugin.ExecOption) error,
		) error {
			// Register Service
			reg := service.NewRegistry()
			reg.AddServices(log.WithFields(log.Fields{
				"plugin":  plug.Name,
				"command": path.Join(c.Path...),
			}))
			reg.AddServices(reportService)

			// Next Plugin
			return next(ctx, reg.Bind())
		}), plugin.WithExecParallelism(t)); err != nil {
		return err
	}
	return nil
}

func init() {

	scanHostCmd.AddCommand(cmd.MapImageCommand(scanHostImageCmd, scanImage))
	scanHostCmd.AddCommand(cmd.MapContainerCommand(scanHostContainerCmd, scanContainer))

	scanHostCmd.PersistentFlags().Int("threads", 5, "threads for scan action")
	scanHostCmd.PersistentFlags().StringP("output", "o", "report.json", "output filepath of report")
	scanHostCmd.PersistentFlags().StringP("glob", "g", "", "specifies the pattern of plugin file to find")
}
