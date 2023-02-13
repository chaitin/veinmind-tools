package scan

import (
	"context"
	"net/url"
	"path/filepath"
	"regexp"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/remote"
	"github.com/chaitin/veinmind-common-go/pkg/auth"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/log"
	"github.com/distribution/distribution/reference"
	"github.com/gogf/gf/text/gstr"
	"github.com/pkg/errors"

	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/registry"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/target"
)

func Registry(ctx context.Context, t *target.Target, runtime api.Runtime) error {
	remoteRuntime, ok := runtime.(*remote.Runtime)
	if !ok {
		return errors.New("unexpect remote runtime")
	}

	// check target registry domain
	var (
		domain   string
		username string
		password string
		find     bool
	)

	// try url parse
	parsed, err := url.Parse(t.Value)
	if err == nil {
		if parsed.Scheme == "http" || parsed.Scheme == "https" {
			domain = parsed.Host
		}
	}

	// use raw value
	if domain == "" {
		domain = t.Value
	}

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

	// fetch catalog
	repos, err := client.GetRepos(ctx)
	if err != nil {
		return errors.Wrapf(err, "can't get repos for registry %s", domain)
	}

	// filter repos
	var filter []string
	if t.Opts.CatalogFilterRegex != "" {
		r := regexp.MustCompile(t.Opts.CatalogFilterRegex)
		for _, repo := range repos {
			if r.MatchString(repo) {
				filter = append(filter, repo)
			}
		}
		repos = filter
	}

	// load remote image
	var loaded []string
	for _, repo := range repos {
		ref, err := reference.ParseNormalizedNamed(repo)
		if err != nil {
			continue
		}

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
