package main

import (
	_ "embed"
	"errors"
	"path/filepath"
	"strings"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/containerd"
	"github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/pkg/auth"
	"github.com/chaitin/veinmind-common-go/registry"
	commonRuntime "github.com/chaitin/veinmind-common-go/runtime"
	"github.com/distribution/distribution/reference"
	"github.com/spf13/cobra"
)

var scanRegistryCmd = &cmd.Command{
	Use:   "scan-registry",
	Short: "perform registry scan command",
}

var scanRegistryImageCmd = &cmd.Command{
	Use:     "image",
	Short:   "perform registry image scan command",
	PreRunE: scanPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			err             error
			c               commonRuntime.Client
			veinmindRuntime api.Runtime
		)

		server, _ := cmd.Flags().GetString("server")
		config, _ := cmd.Flags().GetString("config")
		namespace, _ := cmd.Flags().GetString("namespace")
		runtime, _ := cmd.Flags().GetString("runtime")
		// tags, _ := cmd.Flags().GetStringSlice("tags")

		// fix: no config need not join path
		if config != "" && parallelContainerMode {
			config = filepath.Join(resourceDirectoryPath, config)
		}

		switch runtime {
		case "docker":
			if config == "" {
				c, err = commonRuntime.NewDockerClient()
			} else {
				authConfig, err := auth.ParseAuthConfig(config)
				if err != nil {
					return err
				}
				c, err = commonRuntime.NewDockerClient(commonRuntime.WithAuth(*authConfig))
			}
			if err != nil {
				return err
			}

			veinmindRuntime, err = docker.New()
			if err != nil {
				return err
			}
		case "containerd":
			c, err = commonRuntime.NewContainerdClient()
			if err != nil {
				return err
			}

			veinmindRuntime, err = containerd.New()
			if err != nil {
				return err
			}
		default:
			return errors.New("runtime not match")
		}

		// If no repo is specified, then query all repo through catalog
		repos := []string{}
		if len(args) == 0 {
			r, err := registry.NewClient(registry.WithAuthFromPath(config))
			if err != nil {
				return err
			}

			repos, err = r.GetRepos(server)
			if err != nil {
				return err
			}
		} else {
			// If it doesn't start with registry, autofill registry
			for _, r := range args {
				rParse, err := reference.Parse(r)
				if err != nil {
					log.Error(err)
					continue
				}

				repos = append(repos, rParse.String())
			}
		}

		if namespace != "" {
			namespaceMaps := map[string][]string{}
			for _, repo := range repos {
				rNamed, err := reference.ParseNamed(repo)
				if err != nil {
					log.Error(err)
					continue
				}

				p := reference.Path(rNamed)
				ns := strings.Split(p, "/")[0]
				namespaceMaps[ns] = append(namespaceMaps[ns], repo)
			}

			_, ok := namespaceMaps[namespace]
			if ok {
				repos = namespaceMaps[namespace]
			} else {
				return errors.New("Namespace doesn't match any repos")
			}
		}

		// get repos tags
		reposN := []string{}
		for _, repo := range repos {
			repoR, err := reference.Parse(repo)
			if err != nil {
				reposN = append(reposN, repo)
				continue
			}

			switch repoR.(type) {
			case reference.Tagged:
				reposN = append(reposN, repo)
				continue
			}

			// get repos tags from remote registry
			r, err := registry.NewClient(registry.WithAuthFromPath(config))
			if err != nil {
				return err
			}

			tags, err := r.GetRepoTags(repo)
			if err != nil {
				reposN = append(reposN, repo)
				continue
			}

			for _, tag := range tags {
				reposN = append(reposN, strings.Join([]string{repo, tag}, ":"))
			}
		}
		repos = reposN

		for _, repo := range repos {
			log.Infof("Start pull image: %#v\n", repo)
			r, err := c.Pull(cmd.Context(), repo)
			if err != nil {
				log.Errorf("Pull image error: %#v\n", err.Error())
				continue
			}
			log.Infof("Pull image success: %#v\n", repo)

			var (
				rNamed reference.Named
			)

			switch c.(type) {
			case *commonRuntime.DockerClient:
				rNamed, err = reference.ParseDockerRef(r)
				if err != nil {
					log.Error(err)
					continue
				}

				domain := reference.Domain(rNamed)
				if domain == "index.docker.io" || domain == "docker.io" {
					repo = reference.Path(rNamed)
					if (strings.Split(repo, "/")[0] == "library" || strings.Split(repo, "/")[0] == "_") && len(strings.Split(repo, "/")) >= 2 {
						repo = strings.Join(strings.Split(repo, "/")[1:], "")
					}
				}
			case *commonRuntime.ContainerdClient:
				repo = r
			}

			ids, err := veinmindRuntime.FindImageIDs(repo)
			switch c.(type) {
			case *commonRuntime.DockerClient:
				if len(ids) > 0 {
					for _, id := range ids {
						image, err := veinmindRuntime.OpenImageByID(id)
						if err != nil {
							log.Error(err)
							continue
						}

						err = scanImage(cmd, image)
						if err != nil {
							log.Error(err)
							continue
						}
					}

					for _, id := range ids {
						err = c.Remove(cmd.Context(), id)
						if err != nil {
							log.Error(err)
						} else {
							log.Infof("Remove image success: %#v\n", repo)
						}
					}
				}
			case *commonRuntime.ContainerdClient:
				image, err := veinmindRuntime.OpenImageByID(r)
				if err != nil {
					log.Error(err)
					continue
				}

				var (
					repoRef string
				)
				repoRefs, err := image.RepoRefs()
				if len(repoRefs) > 0 {
					repoRef = repoRefs[0]
				} else {
					repoRef = image.ID()
				}

				err = scanImage(cmd, image)
				if err != nil {
					log.Error(err)
				}

				err = c.Remove(cmd.Context(), repoRef)
				if err != nil {
					log.Error(err)
				} else {
					log.Infof("Remove image success: %#v\n", repo)
				}
			}
		}

		return nil
	},
	PostRunE: scanPostRunE,
}

func init() {
	scanRegistryCmd.AddCommand(scanRegistryImageCmd)

	scanRegistryCmd.PersistentFlags().Int("threads", 5, "threads for scan action")
	scanRegistryCmd.PersistentFlags().StringP("output", "o", "report.json", "output filepath of report")
	scanRegistryCmd.PersistentFlags().StringP("glob", "g", "", "specifies the pattern of plugin file to find")

	scanRegistryImageCmd.Flags().StringP("runtime", "r", "docker", "specifies the runtime of registry client to use")
	scanRegistryImageCmd.Flags().StringP("server", "s", "index.docker.io", "server address of registry")
	scanRegistryImageCmd.Flags().StringP("config", "c", "", "auth config path")
	scanRegistryImageCmd.Flags().StringP("namespace", "n", "", "namespace of repo")
	scanRegistryImageCmd.Flags().StringSliceP("tags", "t", []string{"latest"}, "tags of repo")
}
