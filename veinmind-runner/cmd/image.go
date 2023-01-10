package main

import (
	"context"
	"errors"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/containerd"
	"github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/libveinmind/go/plugin/service"
	"github.com/chaitin/libveinmind/go/remote"
	"github.com/chaitin/veinmind-common-go/pkg/auth"
	"github.com/distribution/distribution/reference"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	RegistryRemote "github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/rs/xid"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

var scanImageCmd = &cmd.Command{
	Use:      "image",
	Short:    "perform image scan ",
	PreRunE:  scanReportPreRunE,
	RunE:     ScanImage,
	PostRunE: scanReportPostRunE,
}

func ScanImage(c *cmd.Command, args []string) error {
	wg := &sync.WaitGroup{}
	wg.Add(len(args))
	for _, value := range args {
		handler := ScanImageParser(value)
		go func(handler Handler, value string, wg *sync.WaitGroup) {
			defer wg.Done()
			err := handler(c, value)
			if err != nil {
				log.Error(err)
			}
		}(handler, value, wg)
	}
	wg.Wait()
	return nil
}

func ScanImageDocker(c *cmd.Command, arg string) error {
	regex := "docker?:?(.*)"
	compileRegex := regexp.MustCompile(regex)
	matchArr := compileRegex.FindStringSubmatch(arg)
	r, err := docker.New()
	if err != nil {
		return err
	}
	ids, err := r.FindImageIDs(matchArr[1])

	for _, id := range ids {
		image, err := r.OpenImageByID(id)
		if err != nil {
			log.Error(err)
			return err
		}
		scan(c, image)
	}
	return nil
}

func ScanImageContainerd(c *cmd.Command, arg string) error {
	regex := "containerd?:?(.*)"
	compileRegex := regexp.MustCompile(regex)
	matchArr := compileRegex.FindStringSubmatch(arg)
	r, err := containerd.New()
	if err != nil {
		return err
	}
	ids, err := r.FindImageIDs(matchArr[1])
	for _, id := range ids {
		image, err := r.OpenImageByID(id)
		if err != nil {
			log.Error(err)
			return err
		}
		scan(c, image)
	}
	return nil
}

func ScanImageRegistry(c *cmd.Command, arg string) error {
	var (
		err                error
		username, password string
		remoteRuntime      *remote.Runtime
		errAssert          bool
	)
	ids := make([]string, 0)
	paths := make(map[string]string, 0)
	authConfig := &auth.AuthConfig{
		Auths: nil,
	}
	regex := "registry?:?(.*)"
	compileRegex := regexp.MustCompile(regex)
	matchArr := compileRegex.FindStringSubmatch(arg)
	registryString := matchArr[1]
	parserRegistry, repos := RegistryParser(registryString, RegistryRemote.WithAuth(&authn.Basic{
		Username: username,
		Password: password,
	}))
	if err != nil {
		return err
	}
	config, err := c.Flags().GetString("config")
	if err != nil {
		return err
	}
	if config != "" {
		if parallelContainerMode {
			config = filepath.Join(resourceDirectoryPath, config)
		}
		authConfig, err = auth.ParseAuthConfig(config)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	//获取registry里的所有的image

	if parserRegistry[0] == "index.docker.io" {
		log.Warnf("found server: docker \nCurrently, docker.io authentication is not supported, so %#v is automatically scanned as a public image without authentication information", registryString)
	} else {
		for _, auth := range authConfig.Auths {
			if strings.Contains(auth.Registry, parserRegistry[0]) {
				username = auth.Username
				password = auth.Password
			}

		}
	}
	for i := 0; i < len(repos); i++ {
		repos[i] = parserRegistry[0] + "/" + repos[i]
	}

	//将registry中所有的image Load进来
	for _, repo := range repos {
		path := filepath.Join("/tmp/", xid.NewWithTime(time.Now()).String())
		paths[repo] = path
		runtime, err := remote.New(path)
		if err != nil {
			log.Error(err)
			continue
		}
		log.Infof("Pull image success: %#v\n", repo)
		remoteRuntime, errAssert = runtime.(*remote.Runtime)
		if errAssert != true {
			log.Error(err)
			continue
		}
		_, err = remoteRuntime.Load(repo, remote.WithAuth(username, password))
		if err != nil {
			log.Error(err)
			continue
		}
	}

	//判断是否指定了image名称,如果没指定就根据输入解析namespace,提取指定namespace下的所有image
	if parserRegistry[2] != "" {
		tmp, err := remoteRuntime.FindImageIDs(registryString)
		for _, id := range tmp {
			ids = append(ids, id)
		}
		if err != nil {
			log.Error(err)
			return err
		}
	} else {
		//解析namespace
		if parserRegistry[1] != "" {
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

			_, ok := namespaceMaps[parserRegistry[1]]
			if ok {
				repos = namespaceMaps[parserRegistry[1]]
			} else {
				return errors.New("Namespace doesn't match any repos")
			}
		}
		for _, repo := range repos {
			tmp, err := remoteRuntime.FindImageIDs(repo)
			for _, id := range tmp {
				ids = append(ids, id)
			}
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}

	//扫描
	for _, id := range ids {
		image, err := remoteRuntime.OpenImageByID(id)
		if err != nil {
			log.Error(err)
			continue
		}
		repoRef, err := image.RepoRefs()
		if err != nil {
			log.Error(err)
			continue
		}
		scan(c, image)
		if err != nil {
			log.Error(err)
			continue
		}
		defer func() {
			for _, ref := range repoRef {
				_, errClose := os.Stat(paths[ref])
				if errClose == nil {
					errRemove := os.RemoveAll(paths[ref])
					if errRemove != nil {
						log.Error(errRemove)
					}
					log.Infof("Remove image success: %#v\n", paths[ref])
				} else {
					log.Error(errClose)
				}
			}
		}()
	}
	return nil
}

func scan(c *cmd.Command, image api.Image) error {
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
			opts := make([]plugin.ExecOption, 0)
			opts = append(opts, reg.Bind())
			if value, ok := pluginArgsMap[plug.Name]; ok == true {
				opts = append(opts, plugin.WithPrependArgs(value...))
			}
			reg.AddServices(log.WithFields(log.Fields{
				"plugin":  plug.Name,
				"command": path.Join(c.Path...),
			}))
			reg.AddServices(reportService)

			// Next Plugin
			return next(ctx, opts...)
		}), plugin.WithExecParallelism(t)); err != nil {
		return err
	}
	return nil
}

func RegistryParser(arg string, auths RegistryRemote.Option) ([]string, []string) {
	res := make([]string, 0)
	repos := make([]string, 0)
	splitRes := strings.Split(arg, "/")
	if len(splitRes) == 3 {
		res = append(res, splitRes[0], splitRes[1], splitRes[2])
	} else if len(splitRes) == 2 {
		registryAddr, err := name.NewRegistry(splitRes[0])
		if err != nil {
			log.Error(err)
		}
		repos, err = RegistryRemote.Catalog(context.Background(), registryAddr, auths)
		if err == nil {
			res = append(res, splitRes[0], splitRes[1], "")
		} else {
			res = append(res, "index.docker.io", splitRes[0], splitRes[1])
		}
	} else if len(splitRes) == 1 {
		res = append(res, "index.docker.io", "", arg)
	}
	registryAddr, err := name.NewRegistry(res[0])
	if err != nil {
		log.Error(err)
		return res, nil
	}
	if res[0] == "index.docker.io" {
		log.Warnf("found server: docker \nCurrently, docker.io authentication is not supported, so it is automatically scanned as a public image without authentication information")
		auths = RegistryRemote.WithAuth(&authn.Basic{
			Username: "",
			Password: "",
		})
	}
	repos, err = RegistryRemote.Catalog(context.Background(), registryAddr, auths)
	if err != nil {
		log.Error(err)
		return res, nil
	}
	return res, repos
}

func ScanImageParser(arg string) Handler {
	regex := "(docker|containerd|tarball|registry):(.*)"
	compileRegex := regexp.MustCompile(regex)
	matchArr := compileRegex.FindStringSubmatch(arg)
	switch matchArr[1] {
	case DOCKER:
		return ScanImageDocker
	case CONTAINERD:
		return ScanImageContainerd
	case REGISTRY:
		return ScanImageRegistry
	default:
		log.Errorf("please input right args, available: docker,containerd,registry")
	}
	return nil
}
