package authz

import (
	"context"

	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/authz/route"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/reporter"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/scan"
)

func HandleWebhookImagePush(ctx context.Context, policy Policy, postData route.PullandPushData) (chan []reporter.ReportEvent, error) {
	eventListCh := make(chan []reporter.ReportEvent, 1)
	if postData.Operator == "webhook" || postData.Type != "PUSH_ARTIFACT" {
		return nil, nil
	}
	imageNames, err := route.GetImageNames(postData)
	if err != nil {
		return nil, err
	}
	var result []reporter.ReportEvent
	for _, img := range imageNames {
		report, err := scan.ScanLocalImage(ctx, img,
			policy.EnabledPlugins, policy.PluginParams)
		if err != nil {
			log.Error(err)
			continue
		}
		result = append(result, report...)
	}
	eventListCh <- result
	return eventListCh, nil
}
