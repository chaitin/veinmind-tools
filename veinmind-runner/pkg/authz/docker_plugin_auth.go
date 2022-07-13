package authz

import (
	"context"
	"fmt"
	"github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/authz/action"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/authz/route"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/reporter"
	scankit "github.com/chaitin/veinmind-tools/veinmind-runner/pkg/scan"
	"github.com/docker/docker/pkg/authorization"
	"time"
)

type DockerPluginHandler func(policy Policy, req *authorization.Request) (bool, error)

func handleContainerCreate(policy Policy, req *authorization.Request, runnerReporter *reporter.Reporter, reportService *report.ReportService) (bool, error) {
	imageName, err := route.GetImageNameFromBodyParam(req.RequestURI, req.RequestHeaders["Content-Type"], "Image", req.RequestBody)
	if err != nil {
		return true, err
	}

	err = scankit.ScanLocalImage(context.Background(), imageName, policy.EnabledPlugins, policy.PluginParams, reportService)
	if err != nil {
		log.Error(err)
	}

	riskLevelFilter := make(map[string]struct{})
	for _, level := range policy.RiskLevelFilter {
		riskLevelFilter[level] = struct{}{}
	}
	events, _ := runnerReporter.GetEvents()
	for _, event := range events {
		if _, ok := riskLevelFilter[toLevelStr(event.Level)]; !ok {
			continue
		}

		if policy.Block {
			return false, nil
		}
	}

	return true, nil
}

var imageCreateMap = action.NewMap()

func handleImageCreate(policy Policy, req *authorization.Request, runnerReporter *reporter.Reporter, reportService *report.ReportService) (bool, error) {
	imageName, err := route.GetImageNameFromUrlParam(req.RequestURI, "fromImage")
	if err != nil {
		return true, err
	}

	imageActionId := fmt.Sprintf("%s-%d", imageName, time.Now().UnixMicro())
	if imageCreateMap.Count(imageName) > 1 {
		return false, nil
	}

	imageCreateMap.Store(imageActionId)
	go func() {
		defer imageCreateMap.Delete(imageActionId)

		ticker := time.NewTicker(time.Second * 5)
		runtime, _ := docker.New()
		for {
			select {
			case <-time.After(time.Minute * 10):
				return
			case <-ticker.C:
				imageIds, err := runtime.FindImageIDs(imageName)
				if err != nil || len(imageIds) < 1 {
					break
				}

				err = scankit.ScanLocalImage(context.Background(), imageName, policy.EnabledPlugins, policy.PluginParams, reportService)
				if err != nil {
					log.Error(err)
				}
				return
			}
		}
	}()

	return true, nil
}

func handleImagePush(policy Policy, req *authorization.Request, runnerReporter *reporter.Reporter, reportService *report.ReportService) (bool, error) {
	imageName, err := route.GetImageNameFromUri(req.RequestURI)
	if err != nil {
		return true, err
	}

	err = scankit.ScanLocalImage(context.Background(), imageName, policy.EnabledPlugins, policy.PluginParams, reportService)
	if err != nil {
		log.Error(err)
	}

	riskLevelFilter := make(map[string]struct{})
	for _, level := range policy.RiskLevelFilter {
		riskLevelFilter[level] = struct{}{}
	}
	events, _ := runnerReporter.GetEvents()
	for _, event := range events {
		if _, ok := riskLevelFilter[toLevelStr(event.Level)]; !ok {
			continue
		}

		if policy.Block {
			return false, nil
		}
	}

	return true, nil
}
