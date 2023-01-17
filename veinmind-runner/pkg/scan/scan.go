package scan

import (
	"context"
	"path"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/libveinmind/go/plugin/service"
	"github.com/chaitin/libveinmind/go/plugin/specflags"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/reporter"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/target"
)

// DispatchTask declare func that how to scan a target object
type DispatchTask func(ctx context.Context, targets []*target.Target) error

func FindTargetPlugins(ctx context.Context, enablePlugins []string) ([]*plugin.Plugin, error) {
	ps, err := plugin.DiscoverPlugins(ctx, ".")
	if err != nil {
		return nil, err
	}
	pluginMap := make(map[string]*plugin.Plugin)
	for _, p := range ps {
		pluginMap[p.Name] = p
	}
	// find the intersection of plugins
	// between found in runner and user specified
	finalPs := []*plugin.Plugin{}
	for _, item := range enablePlugins {
		if p, ok := pluginMap[item]; ok {
			finalPs = append(finalPs, p)
		}
	}
	return finalPs, nil
}

func ScanLocalImage(ctx context.Context, imageName string,
	enabledPlugins []string, pluginParams []string) (events []reporter.ReportEvent, err error) {
	reportService := report.NewReportService()
	runnerReporter, _ := reporter.NewReporter()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go startReportService(ctx, runnerReporter, reportService)
	go runnerReporter.Listen()
	defer runnerReporter.StopListen()

	veinmindRuntime, err := docker.New()
	if err != nil {
		return events, err
	}
	imageIDs, err := veinmindRuntime.FindImageIDs(imageName)
	if err != nil {
		return events, err
	}
	finalPs, err := FindTargetPlugins(ctx, enabledPlugins)
	if err != nil {
		return events, err
	}
	for _, id := range imageIDs {
		image, err := veinmindRuntime.OpenImageByID(id)
		if err != nil {
			log.Error(err)
			continue
		}
		err = ScanImage(ctx, finalPs, image, reportService,
			specflags.WithSpecFlags(pluginParams))
		if err != nil {
			log.Error(err)
		}
	}
	return runnerReporter.GetEvents()
}

func ScanImage(ctx context.Context, rang plugin.ExecRange, image api.Image,
	reportService *report.ReportService, opts ...plugin.ExecOption) error {
	opts = append(opts, plugin.WithExecInterceptor(func(
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
	}))
	return cmd.ScanImage(ctx, rang, image, opts...)
}

func startReportService(ctx context.Context,
	runnerReporter *reporter.Reporter, reportService *report.ReportService) {
	for {
		select {
		case <-ctx.Done():
			return
		case evt := <-reportService.EventChannel:
			runnerReporter.EventChannel <- evt
		}
	}
}
