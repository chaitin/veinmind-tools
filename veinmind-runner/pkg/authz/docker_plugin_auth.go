package authz

import (
	"context"
	"fmt"
	"time"

	"github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/authz/route"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/reporter"
	scankit "github.com/chaitin/veinmind-tools/veinmind-runner/pkg/scan"
	"github.com/docker/docker/pkg/authorization"
)

func handleContainerCreate(policy Policy, req *authorization.Request,
	runnerReporter *reporter.Reporter, reportService *report.ReportService) (bool, error) {
	defer runnerReporter.StopListen()

	imageName, err := route.GetImageNameFromBodyParam(req.RequestURI,
		req.RequestHeaders["Content-Type"], "Image", req.RequestBody)
	if err != nil {
		return true, err
	}
	err = scankit.ScanLocalImage(context.Background(), imageName,
		policy.EnabledPlugins, policy.PluginParams, reportService)
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

var imageCreateMap = newHandleMap()

func handleImageCreate(policy Policy, req *authorization.Request, runnerReporter *reporter.Reporter, reportService *report.ReportService) (bool, error) {
	imageName, err := route.GetImageNameFromUrlParam(req.RequestURI, "fromImage")
	if err != nil {
		runnerReporter.StopListen()
		return true, err
	}

	handleId := fmt.Sprintf("%s-%d", imageName, time.Now().UnixMicro())
	if imageCreateMap.Count(handleId) > 1 {
		runnerReporter.StopListen()
		return false, nil
	}

	imageCreateMap.Store(handleId)
	go func() {
		defer func() {
			imageCreateMap.Delete(handleId)
			runnerReporter.StopListen()
		}()

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

				handleReportAlert(policy, runnerReporter)
				return
			}
		}
	}()

	return true, nil
}

func handleImagePush(policy Policy, req *authorization.Request, runnerReporter *reporter.Reporter, reportService *report.ReportService) (bool, error) {
	defer runnerReporter.StopListen()

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
