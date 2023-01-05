package main

import (
	_ "embed"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/libveinmind/go/remote"
	"github.com/chaitin/veinmind-common-go/pkg/auth"
	"github.com/chaitin/veinmind-common-go/registry"
	"github.com/distribution/distribution/reference"
	"github.com/rs/xid"
	"github.com/spf13/cobra"
)

var scanRegistryCmd = &cmd.Command{
	Use:   "scan-registry",
	Short: "perform registry scan command",
}

var scanRegistryImageCmd = &cmd.Command{
	Use:      "image",
	Short:    "perform registry image scan command",
	PreRunE:  scanPreRunE,
	RunE:     ScanRegistry,
	PostRunE: scanPostRunE,
}

func ScanRegistry(cmd *cobra.Command, args []string) error {
	server, _ := cmd.Flags().GetString("server")
	config, _ := cmd.Flags().GetString("config")
	namespace, _ := cmd.Flags().GetString("namespace")
	// tags, _ := cmd.Flags().GetStringSlice("tags")
	var (
		r                  *registry.Client
		err                error
		ids                []string
		username, password string
	)
	authConfig := &auth.AuthConfig{
		Auths: nil,
	}
	// fix: no config need not join path
	if config == "" {
		r, err = registry.NewClient()
		if err != nil {
			return err
		}
	} else {
		if parallelContainerMode {
			config = filepath.Join(resourceDirectoryPath, config)
		}
		authConfig, err = auth.ParseAuthConfig(config)
		if err != nil {
			log.Error(err)
			return err
		}
		r, err = registry.NewClient(registry.WithAuthFromPath(config))
		if err != nil {
			return err
		}

	}

	// If no repo is specified, then query all repo through catalog
	repos := []string{}
	if len(args) == 0 {
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
		path := filepath.Join("/tmp/", xid.NewWithTime(time.Now()).String())
		runtime, err := remote.New(path)
		if err != nil {
			log.Error(err)
			continue
		}
		remoteRuntime, errAssert := runtime.(*remote.Runtime)
		if errAssert != true {
			log.Error(err)
			continue
		}
		for _, value := range authConfig.Auths {
			if strings.HasPrefix(repo, value.Registry) {
				username = value.Username
				password = value.Password
			}
		}
		ids, _ = remoteRuntime.Load(repo, remote.WithAuth(username, password))
		log.Infof("Pull image success: %#v\n", repo)
		defer func() {
			_, errClose := os.Stat(path)
			if errClose == nil {
				errRemove := os.RemoveAll(path)
				if errRemove != nil {
					log.Error(errRemove)
				}
				log.Infof("Remove image success: %#v\n", repo)
			} else {
				log.Error(errClose)
			}
		}()
		for _, id := range ids {
			image, err := runtime.OpenImageByID(id)
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

	}
	return nil
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
