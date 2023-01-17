package target

import (
	"context"
	"path"

	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/libveinmind/go/plugin/service"
	"github.com/chaitin/libveinmind/go/plugin/specflags"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/plugind"
)

type Target struct {
	Protol         Protol
	Value          string
	Opts           *Options
	Plugins        []*plugin.Plugin
	ServiceManager *plugind.Manager
	ReportService  *report.ReportService
}

// WithDefaultOptions Transfer Target.Options to plugin.ExecOption
func (t *Target) WithDefaultOptions(opts ...Option) []plugin.ExecOption {
	for _, o := range opts {
		o(t.Opts)
	}
	// Exec Options
	popt := make([]plugin.ExecOption, 0)

	// plugin's params Interceptor
	if len(t.Opts.SpecFlags) > 0 {
		popt = append(popt, specflags.WithSpecFlags(t.Opts.SpecFlags))
	}
	// service add Interceptor
	popt = append(popt, plugin.WithExecInterceptor(func(
		ctx context.Context, plug *plugin.Plugin, c *plugin.Command,
		next func(context.Context, ...plugin.ExecOption) error,
	) error {
		// Init Service
		log.Infof("Discovered plugin: %#v\n", plug.Name)
		// IaC need not init any service
		err := t.ServiceManager.StartWithContext(ctx, plug.Name)
		if err != nil {
			log.Errorf("%#v can not work: %#v\n", plug.Name, err)
			return err
		}
		// Register Service
		reg := service.NewRegistry()
		reg.AddServices(log.WithFields(log.Fields{
			"plugin":  plug.Name,
			"command": path.Join(c.Path...),
		}))
		reg.AddServices(t.ReportService)

		// Next Plugin
		return next(ctx, reg.Bind())
	}))
	// limit
	if t.Opts.Thread != 0 {
		popt = append(popt, plugin.WithExecParallelism(t.Opts.Thread))
	}

	return popt
}

// NewTargets parse arg to Target
// docker: || "" 			   -> scan all with docker runtime
// docker:imageRef || imageRef -> scan imageRef with docker runtime
// containerd:  			   -> scan all with containerd runtime
// containerd:imageRef 		   -> scan imageRef with containerd runtime
// registry: 				   -> scan all with remote runtime(do not support docker.io)
// registry: imageRef		   -> scan imageRef with remote runtime
func NewTargets(cmd *cmd.Command, args []string, plugins []*plugin.Plugin, serviceManager *plugind.Manager, reportService *report.ReportService, opts ...Option) []*Target {
	targets := make([]*Target, 0)
	// extends flags
	objArgs, specFlags := splitArgs(cmd, args)

	if len(specFlags) > 0 {
		opts = append(opts, WithSpecFlags(specFlags))
	}

	options := &Options{}
	for _, o := range opts {
		o(options)
	}

	if len(objArgs) == 0 {
		args = append(objArgs, "")
	}

	for _, arg := range objArgs {
		protol, value := ParseProto(cmd.Use, arg)
		if protol == UNKNOWN {
			log.Warnf("Identified proto with arg: %s", arg)
			continue
		}
		targets = append(targets, &Target{
			Protol:         protol,
			Value:          value,
			Opts:           options,
			Plugins:        plugins,
			ReportService:  reportService,
			ServiceManager: serviceManager,
		})
	}
	return targets
}

func splitArgs(cmd *cmd.Command, args []string) ([]string, []string) {
	if cmd.ArgsLenAtDash() >= 0 {
		return args[:cmd.ArgsLenAtDash()], args[cmd.ArgsLenAtDash():]
	}
	return args, []string{}
}
