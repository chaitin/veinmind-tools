package scan

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/containerd"
	"github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/libveinmind/go/remote"
	"github.com/chaitin/veinmind-common-go/pkg/auth"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/registry"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/target"
	"github.com/gogf/gf/text/gstr"
	"github.com/rs/xid"
	"golang.org/x/sync/errgroup"
)

func DispatchImages(ctx context.Context, targets []*target.Target) error {
	errG := errgroup.Group{}
	for _, obj := range targets {
		errG.Go(func() error {
			switch obj.Protol {
			case target.DOCKERD:
				r, err := docker.New()
				if err != nil {
					return err
				}
				return HostImage(ctx, obj, r)
			case target.CONTAINERD:
				r, err := containerd.New()
				if err != nil {
					return err
				}
				return HostImage(ctx, obj, r)
			case target.REGISTRY:
				path := filepath.Join(obj.Opts.TempPath, xid.NewWithTime(time.Now()).String())
				r, err := remote.New(path)
				if err != nil {
					return err
				}
				return RegistryImage(ctx, obj, r)
			default:
				return errors.New(fmt.Sprintf("individual image protol: %s", obj.Protol))
			}
		})
	}
	return errG.Wait()
}

func HostImage(ctx context.Context, t *target.Target, runtime api.Runtime) error {
	var ids []string
	var err error
	// scanAll
	if t.Value == "" {
		ids, err = runtime.ListImageIDs()
	} else {
		ids, err = runtime.FindImageIDs(t.Value)
	}
	if err != nil {
		return err
	}
	for _, id := range ids {
		image, err := runtime.OpenImageByID(id)
		if err != nil {
			return err
		}
		if err := doImage(ctx, t.Plugins, image, t.WithDefaultOptions()...); err != nil {
			log.Errorf("scan Image %s error : %s", image.ID(), err)
		}
	}
	return nil
}

func RegistryImage(ctx context.Context, t *target.Target, runtime api.Runtime) error {

	remoteRuntime, ok := runtime.(*remote.Runtime)
	if !ok {
		return errors.New("unexpect remote runtime")
	}

	var username = ""
	var password = ""

	registryOpt := make([]registry.Option, 0)
	registryOpt = append(registryOpt, registry.WithInsecure(t.Opts.Insecure))

	if t.Opts.ConfigPath != "" {
		config := t.Opts.ConfigPath
		if t.Opts.ParallelContainerMode {
			config = filepath.Join(t.Opts.ResourcePath, t.Opts.ConfigPath)
		}
		authConfig, err := auth.ParseAuthConfig(config)
		if err != nil {
			log.Error("load remote auth config error")
		}
		for _, auth := range authConfig.Auths {
			if gstr.HasPrefix(t.Value, auth.Registry) {
				username = auth.Username
				password = auth.Password
				break
			}
		}
	}

	registryOpt = append(registryOpt, registry.WithLoginByPassword(username, password))
	client, err := registry.NewV2Client(t.Value, registryOpt...)

	if err != nil {
		return err
	}
	repos, err := client.GetRepos(ctx)

	if err != nil {
		return err
	}

	//将registry中所有的image Load进来
	for _, repo := range repos {
		tags, err := client.GetRepoTags(ctx, repo)
		if err != nil {
			log.Error(err)
			continue
		}
		for _, tag := range tags {
			name := strings.Join([]string{repo, tag}, ":")
			_, err = remoteRuntime.Load(name, remote.WithAuth(username, password))
			if err != nil {
				log.Error(err)
				continue
			}
			log.Infof("Load image success: %#v\n", name)
		}
	}

	return HostImage(ctx, t, remoteRuntime)
}

func doImage(ctx context.Context, rang plugin.ExecRange, image api.Image, pluginOpts ...plugin.ExecOption) error {
	refs, err := image.RepoRefs()
	ref := ""
	if err == nil && len(refs) > 0 {
		ref = refs[0]
	} else {
		ref = image.ID()
	}
	log.Infof("Scan image: %#v\n", ref)
	return cmd.ScanImage(ctx, rang, image, pluginOpts...)
}
