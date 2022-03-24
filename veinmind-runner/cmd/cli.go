package main

import (
	"context"
	"errors"
	"github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/libveinmind/go/plugin/service"
	"github.com/chaitin/veinmind-tools/veinmind-common/go/service/report"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/registry"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var (
	ps            []*plugin.Plugin
	ctx           context.Context
	reportService = func() *report.ReportService {
		return report.NewReportService()
	}()
	scanPreRunE = func(c *cobra.Command, args []string) error {
		// Discover Plugins
		ctx = c.Context()
		glob, err := c.Flags().GetString("glob")
		if err == nil && glob != "" {
			ps, err = plugin.DiscoverPlugins(ctx, ".", plugin.WithGlob(glob))
		} else {
			ps, err = plugin.DiscoverPlugins(ctx, ".")
		}
		if err != nil {
			return err
		}
		for _, p := range ps {
			log.Infof("Discovered plugin: %#v\n", p.Name)
		}

		// Event Channel Listen
		go func() {
			for {
				select {
				case evt := <-reportService.EventChannel:
					log.Info(evt)
				}
			}
		}()

		return nil
	}
)

var rootCmd = &cmd.Command{}
var scanHostCmd = &cmd.Command{
	Use:     "scan-host",
	Short:   "perform hosted scan command",
	PreRunE: scanPreRunE,
}
var scanRegistryCmd = &cmd.Command{
	Use:     "scan-registry",
	Short:   "perform registry scan command",
	PreRunE: scanPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, _ := cmd.Flags().GetString("address")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		namespace, _ := cmd.Flags().GetString("namespace")
		reponames, _ := cmd.Flags().GetStringSlice("reponames")
		tags, _ := cmd.Flags().GetStringSlice("tags")

		auth := &registry.Auth{}
		if username != "" && password != "" {
			auth.Username = username
			auth.Password = password
		} else {
			auth = nil
		}

		client, err := registry.NewRegistryClient(address, auth)
		if err != nil {
			return err
		}

		// If no repo is specified, then query all repo through catalog
		repos := []string{}
		if len(reponames) == 0 {
			repos, err = client.GetRepos()
			if err != nil {
				return err
			}
		} else {
			// If it doesn't start with registry, autofill registry
			for _, r := range reponames {
				rSplit := strings.Split(r, "/")
				rNew := ""
				if !strings.EqualFold(rSplit[0], address) {
					rSplitNew := []string{address}
					rSplitNew = append(rSplitNew, rSplit...)
					rNew = strings.Join(rSplitNew, "/")
				} else {
					rNew = r
				}
				repos = append(repos, rNew)
			}
		}

		if namespace != "" {
			namespaceMaps := map[string][]string{}
			for _, repo := range repos {
				repoSplit := strings.Split(repo, "/")
				if len(repoSplit) >= 3 {
					namespace := repoSplit[1]
					namespaceMaps[namespace] = append(namespaceMaps[namespace], repo)
				} else if len(repoSplit) == 2 {
					namespace := repoSplit[0]
					namespaceMaps[namespace] = append(namespaceMaps[namespace], repo)
				} else if len(repoSplit) == 1 {
					namespaceMaps["_"] = append(namespaceMaps["_"], repo)
				}
			}

			_, ok := namespaceMaps[namespace]
			if ok {
				repos = namespaceMaps[namespace]
			} else {
				return errors.New("Namespace doesn't match any repos")
			}
		}

		if len(tags) > 0 {
			reposTemp := []string{}
			for _, repo := range repos {
				rtags, err := client.GetRepoTags(repo)
				if err != nil {
					log.Error(err)
					continue
				}

				for _, t1 := range rtags {
					for _, t2 := range tags {
						if strings.EqualFold(t1, t2) {
							repoSplit := strings.Split(repo, ":")
							if len(repoSplit) == 1 {
								repoSplit = append(repoSplit, t1)
								repoWithTag := strings.Join(repoSplit, ":")
								reposTemp = append(reposTemp, repoWithTag)
							}
						}
					}
				}
			}
			repos = reposTemp
		}

		d, err := docker.New()
		if err != nil {
			return err
		}
		defer func() {
			d.Close()
		}()

		for _, repo := range repos {
			log.Infof("Start pull image: %#v\n", repo)
			r, err := client.Pull(repo)
			if err != nil {
				log.Errorf("Pull image error: %#v\n", err.Error())
				continue
			}

			_, err = ioutil.ReadAll(r)
			if err != nil {
				log.Errorf("Pull image error: %#v\n", err.Error())
				continue
			}
			log.Infof("Pull image success: %#v\n", repo)
			ids, err := d.FindImageIDs(repo)
			defer func() {
				for _, id := range ids {
					_, err := client.Remove(id)
					if err != nil {
						log.Error(err)
					}
					log.Infof("Remove image success: %#v\n", repo)
				}
			}()

			if len(ids) > 0 {
				for _, id := range ids {
					image, err := d.OpenImageByID(id)
					if err != nil {
						log.Error(err)
						continue
					}

					err = scan(cmd, image)
					if err != nil {
						log.Error(err)
						continue
					}
				}
			}
		}

		return nil
	},
}

func scan(c *cmd.Command, image api.Image) error {
	refs, err := image.RepoRefs()
	ref := ""
	if err == nil && len(refs) > 0 {
		ref = refs[0]
	} else {
		ref = image.ID()
	}
	log.Infof("Scan image: %#v\n", ref)
	if err := cmd.ScanImages(ctx, ps, []api.Image{image},
		plugin.WithExecInterceptor(func(
			ctx context.Context, plug *plugin.Plugin, c *plugin.Command,
			next func(context.Context, ...plugin.ExecOption) error,
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
		})); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(cmd.MapImageCommand(scanHostCmd, scan))
	rootCmd.AddCommand(scanRegistryCmd)
	scanHostCmd.Flags().StringP("glob", "g", "", "specifies the pattern of plugin file to find")
	scanRegistryCmd.Flags().StringP("address", "a", "index.docker.io", "server address of registry")
	scanRegistryCmd.Flags().StringP("username", "u", "", "username of registry")
	scanRegistryCmd.Flags().StringP("password", "p", "", "password of registry")
	scanRegistryCmd.Flags().StringP("namespace", "n", "", "namespace of repo")
	scanRegistryCmd.Flags().StringSliceP("reponames", "r", []string{}, "name of repo")
	scanRegistryCmd.Flags().StringSliceP("tags", "t", []string{"latest"}, "tags of repo")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
