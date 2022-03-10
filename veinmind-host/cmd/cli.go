package main

import (
	"context"
	"github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/containerd"
	"github.com/chaitin/libveinmind/go/docker"
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
	RunE: func(c *cmd.Command, args []string) error {
		// Discover Plugins
		ctx := c.Context()
		log.Info("Start discovering")
		ps, err := plugin.DiscoverPlugins(ctx, ".")
		if err != nil {
			return err
		}
		for _, p := range ps {
			log.Infof("Discovered plugin: %#v\n", p.Name)
		}
		log.Info("Done discovering")

		var images []api.Image
		defer func() {
			for _, image := range images {
				_ = image.Close()
			}
		}()

		useDocker, err := c.Flags().GetBool("docker")
		if err == nil && useDocker {
			// Init Docker Runtime
			d, err := docker.New()
			if err != nil {
				return err
			}
			defer func() { _ = d.Close() }()

			imageIDs := []string{}
			if len(args) > 0 {
				for _, arg := range args {
					imageIDsTemp, err := d.FindImageIDs(arg)
					if err != nil {
						continue
					}
					imageIDs = append(imageIDs, imageIDsTemp...)
				}
			} else {
				imageIDs, err = d.ListImageIDs()
				if err != nil {
					return err
				}
			}

			// Get Docker Image
			for _, imageID := range imageIDs {
				image, err := d.OpenImageByID(imageID)
				if err != nil {
					return err
				}
				refs, err := image.RepoRefs()
				if err != nil {
					return err
				}
				log.Infof("Enqueuing image: %#v\n", refs)
				images = append(images, image)
			}
		}

		useContainerd, err := c.Flags().GetBool("containerd")
		if err == nil && useContainerd {
			// Init Containerd Runtime
			cd, err := containerd.New()
			if err != nil {
				return err
			}
			defer func() { _ = cd.Close() }()

			imageIDs := []string{}
			if len(args) > 0 {
				for _, arg := range args {
					imageIDsTemp, err := cd.FindImageIDs(arg)
					if err != nil {
						continue
					}
					imageIDs = append(imageIDs, imageIDsTemp...)
				}
			} else {
				imageIDs, err = cd.ListImageIDs()
				if err != nil {
					return err
				}
			}

			// Get Containerd Image
			for _, imageID := range imageIDs {
				image, err := cd.OpenImageByID(imageID)
				if err != nil {
					return err
				}
				refs, err := image.RepoRefs()
				if err != nil {
					return err
				}
				log.Infof("Enqueuing image: %#v\n", refs)
				images = append(images, image)
			}
		}

		log.Infof("Start scanning")
		if err := cmd.ScanImages(ctx, ps, images,
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
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().Bool("docker", true, "specify \"docker\" as the mode in use")
	scanCmd.Flags().Bool("containerd", true, "specify \"containerd\" as the mode in use")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
