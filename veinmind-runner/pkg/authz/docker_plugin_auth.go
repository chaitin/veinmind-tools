package authz

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/authz/route"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/reporter"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/scan"
	"github.com/distribution/distribution/reference"
	"github.com/docker/docker/pkg/authorization"
)

func handleContainerCreate(policy Policy, req *authorization.Request) (<-chan []reporter.ReportEvent, bool, error) {
	eventListCh := make(chan []reporter.ReportEvent, 1)
	defer close(eventListCh)
	imageName, err := route.GetImageNameFromBodyParam(req.RequestURI, req.RequestHeaders["Content-Type"], "Image", req.RequestBody)
	if err != nil {
		return eventListCh, true, err
	}

	events, err := scan.ScanLocalImage(context.Background(), imageName,
		policy.EnabledPlugins, policy.PluginParams)
	if err != nil {
		log.Error(err)
	}
	eventListCh <- events

	return eventListCh, handlePolicyCheck(policy, events), nil
}

var imageCreateMap sync.Map

func handleImageCreate(policy Policy, req *authorization.Request) (<-chan []reporter.ReportEvent, bool, error) {
	eventListCh := make(chan []reporter.ReportEvent, 1)
	imageName, err := route.GetImageNameFromUrlParam(req.RequestURI, "fromImage")
	if err != nil {
		close(eventListCh)
		return eventListCh, true, err
	}

	_, err = reference.Parse(imageName)
	if err != nil {
		close(eventListCh)
		return eventListCh, true, err
	}

	count := 0
	imageCreateMap.Range(func(key, _ interface{}) bool {
		if strings.HasPrefix(key.(string), imageName) {
			count += 1
		}
		return true
	})
	if count > 1 {
		close(eventListCh)
		return eventListCh, false, nil
	}

	handleId := fmt.Sprintf("%s-%d", imageName, time.Now().UnixMicro())
	imageCreateMap.Store(handleId, struct{}{})
	go func() {
		defer func() {
			imageCreateMap.Delete(handleId)
			close(eventListCh)
		}()

		ticker := time.NewTicker(time.Second * 5)
		runtime, _ := docker.New()
		for {
			select {
			case <-time.After(time.Minute * 10):
				return
			case <-ticker.C:
				imageIds, err := runtime.FindImageIDs(imageName)
				if err != nil {
					log.Error(err)
					break
				}

				if len(imageIds) < 1 {
					break
				}

				events, err := scan.ScanLocalImage(context.Background(), imageName,
					policy.EnabledPlugins, policy.PluginParams)
				if err != nil {
					log.Error(err)
				}

				eventListCh <- events
				return
			}
		}
	}()

	return eventListCh, true, nil
}

func handleImagePush(policy Policy, req *authorization.Request) (<-chan []reporter.ReportEvent, bool, error) {
	eventListCh := make(chan []reporter.ReportEvent, 1)
	defer close(eventListCh)

	var events []reporter.ReportEvent
	imageName, err := route.GetImageNameFromUri(req.RequestURI)
	if err != nil {
		return eventListCh, true, err
	}

	events, err = scan.ScanLocalImage(context.Background(), imageName,
		policy.EnabledPlugins, policy.PluginParams)
	if err != nil {
		log.Error(err)
	}
	eventListCh <- events
	return eventListCh, handlePolicyCheck(policy, events), nil
}

func handleDefaultAction() (<-chan []reporter.ReportEvent, bool, error) {
	eventListCh := make(chan []reporter.ReportEvent, 1)
	defer close(eventListCh)

	return eventListCh, true, nil
}
