package authz

import (
	"context"
	"fmt"
	"github.com/chaitin/veinmind-common-go/service/report/event"
	"strings"
	"sync"
	"time"

	"github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/distribution/distribution/reference"
	"github.com/docker/docker/pkg/authorization"

	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/authz/route"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/log"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/plugind"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/scan"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/target"
)

func handleContainerCreate(policy Policy, req *authorization.Request) (<-chan []*event.Event, bool, error) {
	eventListCh := make(chan []*event.Event, 1<<8)
	defer close(eventListCh)
	imageName, err := route.GetImageNameFromBodyParam(req.RequestURI, req.RequestHeaders["Content-Type"], "Image", req.RequestBody)
	if err != nil {
		return eventListCh, true, err
	}

	events, err := handleScan(context.Background(), imageName,
		policy.EnabledPlugins, policy.PluginParams)
	if err != nil {
		log.GetModule(log.AuthzModuleKey).Error(err)
	}
	eventListCh <- events

	return eventListCh, handlePolicyCheck(policy, events), nil
}

var imageCreateMap sync.Map

func handleImageCreate(policy Policy, req *authorization.Request) (<-chan []*event.Event, bool, error) {
	eventListCh := make(chan []*event.Event, 1<<8)
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
					log.GetModule(log.AuthzModuleKey).Error(err)
					break
				}

				if len(imageIds) < 1 {
					break
				}

				events, err := handleScan(context.Background(), imageName,
					policy.EnabledPlugins, policy.PluginParams)
				if err != nil {
					log.GetModule(log.AuthzModuleKey).Error(err)
				}

				eventListCh <- events
				return
			}
		}
	}()

	return eventListCh, true, nil
}

func handleImagePush(policy Policy, req *authorization.Request) (<-chan []*event.Event, bool, error) {
	eventListCh := make(chan []*event.Event, 1<<8)
	defer close(eventListCh)

	var events []*event.Event
	imageName, err := route.GetImageNameFromUri(req.RequestURI)
	if err != nil {
		return eventListCh, true, err
	}

	events, err = handleScan(context.Background(), imageName,
		policy.EnabledPlugins, policy.PluginParams)
	if err != nil {
		log.GetModule(log.AuthzModuleKey).Error(err)
	}
	eventListCh <- events
	return eventListCh, handlePolicyCheck(policy, events), nil
}

func handleDefaultAction() (<-chan []*event.Event, bool, error) {
	eventListCh := make(chan []*event.Event, 1)
	defer close(eventListCh)

	return eventListCh, true, nil
}

func findTargetPlugins(ctx context.Context, enablePlugins []string) ([]*plugin.Plugin, error) {
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

func handleScan(ctx context.Context, imageName string, enabledPlugins []string, pluginParams []string) (events []*event.Event, err error) {
	var eventsMutex sync.RWMutex
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	serviceManager, err := plugind.NewManager()
	if err != nil {
		return events, err
	}

	reportService := report.NewService(ctx)
	defer close(reportService.EventPool.EventChannel)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case evt := <-reportService.EventPool.EventChannel:
				eventsMutex.Lock()
				events = append(events, evt)
				eventsMutex.Unlock()
			}
		}
	}()

	finalPs, err := findTargetPlugins(ctx, enabledPlugins)
	if err != nil {
		return events, err
	}

	tg := &target.Target{
		Proto: target.DOCKERD,
		Value: imageName,
		Opts: &target.Options{
			SpecFlags: pluginParams,
		},
		Plugins:        finalPs,
		ServiceManager: serviceManager,
		ReportService:  reportService,
	}

	err = scan.DispatchImages(ctx, []*target.Target{tg})

	return events, err
}
