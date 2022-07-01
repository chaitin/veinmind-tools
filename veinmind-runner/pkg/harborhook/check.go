package harborhook

import (
	"context"
	"fmt"
	"os"

	"github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/libveinmind/go/plugin/specflags"
	"github.com/chaitin/veinmind-tools/veinmind-common/go/service/report"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/reporter"
	scanutil "github.com/chaitin/veinmind-tools/veinmind-runner/pkg/scan"
)

type Authorizer struct {
	cfg WebhookConfig
}

func (t *Authorizer) CheckPull(action string, imageName string, server WebhookServer) (err error) {
	dockerClient := server.dockerClient
	_, err = dockerClient.Pull(imageName)
	if err != nil {
		log.Errorf("Pull image error: %#v\n", err.Error())
		return err
	}
	log.Infof("Pull image success: %#v\n", imageName)

	// imageNames := strings.Split(imageName, "@")
	// imageName = imageNames[len(imageNames)-1]
	policy, ok := t.cfg.PolicysMap()[action]
	if !ok {
		return nil
	}
	fmt.Printf("%s policy => %v\n", action, policy)
	var ctx = context.Background()
	var ps []*plugin.Plugin
	reportService := report.NewReportService()
	runnerReporter, err := reporter.NewReporter()

	// get eventReport from Pliugin
	go runnerReporter.Listen()
	go func() {
		for {
			select {
			case evt := <-reportService.EventChannel:
				runnerReporter.EventChannel <- evt
			}
		}
	}()
	if err != nil {
		return err
	}
	// Create a Runtime instance
	veinmindRuntime, err := docker.New()
	if err != nil {
		return err
	}
	imageIDs, err := veinmindRuntime.FindImageIDs(imageName)
	fmt.Println("iamgeIDs => ", imageIDs)
	if err != nil {
		return err
	}
	ps, err = plugin.DiscoverPlugins(ctx, ".")

	if err != nil {
		return err
	}
	pluginMap := make(map[string]*plugin.Plugin)
	for _, p := range ps {
		pluginMap[p.Name] = p
	}
	// find the intersection of plugins
	// between found in runner and user specified
	finalPs := []*plugin.Plugin{}
	for _, item := range policy.EnabledPlugins {
		if p, ok := pluginMap[item]; ok {
			finalPs = append(finalPs, p)
		}
	}
	for _, id := range imageIDs {
		image, err := veinmindRuntime.OpenImageByID(id)
		if err != nil {
			log.Error(err)
			continue
		}
		err = scanutil.ScanImage(ctx, finalPs, image, reportService,
			specflags.WithSpecFlags(policy.PluginParams))
		if err != nil {
			log.Error(err)
		}
	}
	err = dockerClient.Remove(imageName)
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("Remove image success: %#v\n", imageName)
	}

	events, err := runnerReporter.GetEvents()
	if err != nil {
		return err
	}
	output := t.cfg.Log.ReportLogPath
	reportMap := make(map[report.Level]string)
	for _, r := range policy.RiskLevelFilter {
		reportMap[FromLevel[r]] = r
	}
	for _, event := range events {
		if _, ok := reportMap[event.Level]; ok {
			if policy.Alert {
				log.Warn(fmt.Sprintf("%s exec failed, cause plugins find risks in %s", policy.Action, imageName))
			}
			// if err generated from file opreration
			// authz should ignore the err,and return allow/deny to docker compose
			f, err := os.OpenFile(output, os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				log.Error(err)
			} else {
				if len(events) > 0 {
					err = runnerReporter.Write(f)
					if err != nil {
						log.Error(err)
					}
				}
				f.Close()
			}
		}
	}
	fmt.Println(events)
	return nil
}
