package target

import (
	"context"
	"path"

	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	logService "github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/libveinmind/go/plugin/service"
	"github.com/chaitin/libveinmind/go/plugin/specflags"
	"github.com/chaitin/veinmind-common-go/service/report"

	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/log"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/plugind"
)

type Target struct {
	Proto          Proto
	Value          string
	Opts           *Options
	Plugins        []*plugin.Plugin
	ServiceManager *plugind.Manager
	ReportService  *report.Service
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
		log.GetModule(log.TargetModuleKey).Infof("discovered plugin: %#v\n", plug.Name)
		// IaC need not init any service
		err := t.ServiceManager.StartWithContext(ctx, plug.Name)
		if err != nil {
			return err
		}
		// Register Service
		reg := service.NewRegistry()
		reg.AddServices(logService.WithFields(logService.Fields{
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
func NewTargets(cmd *cmd.Command, args []string, plugins []*plugin.Plugin, serviceManager *plugind.Manager, reportService *report.Service, opts ...Option) []*Target {
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
		objArgs = append(objArgs, "")
	}

	for _, arg := range objArgs {
		proto, value := ParseProto(cmd.Name(), arg)
		if proto == UNKNOWN {
			log.GetModule(log.TargetModuleKey).Warnf("can't identified proto with arg: %s", arg)
			continue
		}
		targets = append(targets, &Target{
			Proto:          proto,
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
