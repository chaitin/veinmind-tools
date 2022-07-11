package authz

import (
	"context"
	"fmt"
	"io"

	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/veinmind-common/go/service/report"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/reporter"
	scanutil "github.com/chaitin/veinmind-tools/veinmind-runner/pkg/scan"
)

type checkType uint32

const (
	REMOTE checkType = iota
	LOCAL
)

func CheckImage(imageName string, policy Policy,
	f io.Writer) (bool, error) {
	// before enter this function we have
	// checked that action is in the policyMap
	var ctx = context.Background()
	reportService := report.NewReportService()
	runnerReporter, err := reporter.NewReporter()
	if err != nil {
		log.Error(err)
		return true, err
	}
	// get eventReport from Pliugin
	go runnerReporter.Listen()
	reportCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		for {
			select {
			case <-reportCtx.Done():
				return
			case evt := <-reportService.EventChannel:
				runnerReporter.EventChannel <- evt
			}
		}
	}()

	err = scanutil.ScanLocalImage(ctx, imageName, policy.EnabledPlugins,
		policy.PluginParams, reportService)

	if err != nil {
		log.Error(err)
		return false, err
	}
	events, _ := runnerReporter.GetEvents()

	// Stop reporter listen
	runnerReporter.StopListen()
	reportLevel := policy.RiskLevelFilter
	reportMap := make(map[string]bool)
	for _, r := range reportLevel {
		r = "\"" + r + "\""
		reportMap[r] = true
	}
	for _, event := range events {
		blevel, err := event.Level.MarshalJSON()
		if err != nil {
			log.Error(err)
		}
		fmt.Println(string(blevel))
		if _, ok := reportMap[string(blevel)]; ok {
			if policy.Alert {
				log.Warn(fmt.Sprintf("Image %s has risks!", imageName))
			}
			err = runnerReporter.Write(f)
			if err != nil {
				log.Error(err)
			}
			if policy.Block {
				return true, nil
			}
		}
	}

	return false, nil
}
