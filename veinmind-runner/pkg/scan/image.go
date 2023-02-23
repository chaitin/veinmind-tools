package scan

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/chaitin/libveinmind/go/tarball"
	"github.com/pkg/errors"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/containerd"
	"github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/remote"
	"github.com/chaitin/veinmind-common-go/pkg/auth"
	"github.com/distribution/distribution/reference"
	"github.com/gogf/gf/text/gstr"
	"github.com/rs/xid"
	"golang.org/x/sync/errgroup"

	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/log"

	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/registry"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/target"
)

func DispatchImages(ctx context.Context, targets []*target.Target) error {
	errG := errgroup.Group{}
	for _, obj := range targets {
		o := obj
		errG.Go(func() error {
			switch o.Proto {
			case target.DOCKERD:
				r, err := docker.New()
				if err != nil {
					return err
				}
				return HostImage(ctx, o, r)
			case target.CONTAINERD:
				r, err := containerd.New()
				if err != nil {
					return err
				}
				return HostImage(ctx, o, r)
			case target.REGISTRY_IMAGE:
				path := filepath.Join(o.Opts.TempPath, xid.NewWithTime(time.Now()).String())
				r, err := remote.New(path)
				if err != nil {
					return err
				}
				return RegistryImage(ctx, o, r)
			case target.REGISTRY:
				path := filepath.Join(o.Opts.TempPath, xid.NewWithTime(time.Now()).String())
				r, err := remote.New(path)
				if err != nil {
					return err
				}
				return Registry(ctx, o, r)
			case target.TARBALL:
				path := filepath.Join(o.Opts.TempPath, xid.NewWithTime(time.Now()).String())
				t, err := tarball.New(tarball.WithRoot(path))
				if err != nil {
					return err
				}
				return TarballImage(ctx, o, t)
			default:
				return errors.New(fmt.Sprintf("individual image protol: %s", obj.Proto))
			}
		})
	}
	return errG.Wait()
}

func HostImage(ctx context.Context, t *target.Target, runtime api.Runtime) error {
	var ids []string
	var err error
	// scanAll
	if t.Value == "" || t.Value == "*" {
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
			log.GetModule(log.ScanModuleKey).Errorf("scan Image %s error : %s", image.ID(), err)
		}
	}
	return nil
}

func RegistryImage(ctx context.Context, t *target.Target, runtime api.Runtime) error {
	remoteRuntime, ok := runtime.(*remote.Runtime)
	if !ok {
		return errors.New("unexpect remote runtime")
	}

	var (
		domain   string
		username string
		password string
		find     bool
	)

	// parse image reference domain
	ref, err := reference.ParseDockerRef(t.Value)
	if err != nil {
		return errors.Wrapf(err, "can't parse domain from reference %s", t.Value)
	}
	domain = reference.Domain(ref)

	// parse domain mapping username and password
	registryOpt := make([]registry.Option, 0)
	registryOpt = append(registryOpt, registry.WithInsecure(t.Opts.Insecure))

	if t.Opts.ConfigPath != "" {
		config := t.Opts.ConfigPath
		if t.Opts.ParallelContainerMode {
			config = filepath.Join(t.Opts.ResourcePath, t.Opts.ConfigPath)
		}
		authConfig, err := auth.ParseAuthConfig(config)
		if err != nil {
			log.GetModule(log.ScanModuleKey).Errorf("load remote auth config error, %+v", err)
		} else {
			for _, authEntry := range authConfig.Auths {
				if gstr.Equal(domain, authEntry.Registry) {
					username = authEntry.Username
					password = authEntry.Password
					find = true
					break
				}
			}
		}
	}
	if find {
		registryOpt = append(registryOpt, registry.WithLoginByPassword(username, password))
	}

	// init registry v2 client
	client, err := registry.NewV2Client(domain, registryOpt...)
	if err != nil {
		return errors.Wrapf(err, "can't init registry v2 client for %s", domain)
	}

	// append reference
	normalized, err := reference.ParseNormalizedNamed(t.Value)
	if err != nil {
		return errors.Wrapf(err, "can't parse reference %s", t.Value)
	}
	refs := make([]reference.Named, 0)
	refs = append(refs, normalized)

	// load remote image
	var loaded []string
	for _, ref := range refs {
		// check reference tag
		if _, ok := ref.(reference.NamedTagged); !ok {
			tags, err := client.GetRepoTags(ctx, ref.String())
			if err != nil {
				log.GetModule(log.ScanModuleKey).Errorf("can't get reference tags for %s, err: %+v", ref.String(), err)
				continue
			}

			for _, tag := range tags {
				complete, err := reference.WithTag(ref, tag)
				images, err := remoteRuntime.Load(complete.String(), remote.WithAuth(username, password))
				if err != nil {
					log.GetModule(log.ScanModuleKey).Errorf("can't load remote image for %s, err: %+v", complete.String(), err)
					continue
				}
				loaded = append(loaded, images...)
				log.GetModule(log.ScanModuleKey).Infof("load image success: %#v\n", complete.String())
			}
		} else {
			images, err := remoteRuntime.Load(ref.String(), remote.WithAuth(username, password))
			if err != nil {
				log.GetModule(log.ScanModuleKey).Errorf("can't load remote image for %s, err: %+v", ref.String(), err)
			} else {
				loaded = append(loaded, images...)
				log.GetModule(log.ScanModuleKey).Infof("load image success: %#v\n", ref.String())
			}
		}
	}

	// open and scan image instance
	var (
		images []api.Image
		uniq   map[string]struct{}
	)
	uniq = make(map[string]struct{})
	for _, id := range loaded {
		if _, ok := uniq[id]; ok {
			continue
		} else {
			uniq[id] = struct{}{}
		}

		instance, err := runtime.OpenImageByID(id)
		if err != nil {
			continue
		}
		images = append(images, instance)
	}

	return doImages(ctx, t.Plugins, images, t.WithDefaultOptions()...)
}

func TarballImage(ctx context.Context, t *target.Target, runtime api.Runtime) error {
	tarballRuntime, ok := runtime.(*tarball.Tarball)
	if !ok {
		return errors.New("scan: runtime type not match for tarball")
	}

	_, err := tarballRuntime.Load(t.Value)
	if err != nil {
		return err
	}

	var images []api.Image
	ids, _ := tarballRuntime.ListImageIDs()
	for _, id := range ids {
		image, err := tarballRuntime.OpenImageByID(id)
		if err != nil {
			log.GetModule(log.ScanModuleKey).Error(err)
			continue
		}

		images = append(images, image)
	}

	return doImages(ctx, t.Plugins, images, t.WithDefaultOptions()...)
}

func doImage(ctx context.Context, rang plugin.ExecRange, image api.Image, pluginOpts ...plugin.ExecOption) error {
	refs, err := image.RepoRefs()
	ref := ""
	if err == nil && len(refs) > 0 {
		ref = refs[0]
	} else {
		ref = image.ID()
	}
	log.GetModule(log.ScanModuleKey).Infof("start scan image: %#v\n", ref)
	return cmd.ScanImage(ctx, rang, image, pluginOpts...)
}

func doImages(ctx context.Context, rang plugin.ExecRange, images []api.Image, pluginOpts ...plugin.ExecOption) error {
	for _, image := range images {
		refs, err := image.RepoRefs()
		ref := ""
		if err == nil && len(refs) > 0 {
			ref = refs[0]
		} else {
			ref = image.ID()
		}
		log.GetModule(log.ScanModuleKey).Infof("start scan image: %#v\n", ref)
	}
	return cmd.ScanImages(ctx, rang, images, pluginOpts...)
}
