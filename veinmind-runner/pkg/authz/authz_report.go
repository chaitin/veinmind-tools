package authz

import (
	"context"
	"fmt"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/reporter"
	"io"
)

func toLevelStr(level report.Level) string {
	switch level {
	case report.Low:
		return "Low"
	case report.Medium:
		return "Medium"
	case report.High:
		return "High"
	case report.Critical:
		return "Critical"
	}

	return "None"
}

func startReportService(ctx context.Context, runnerReporter *reporter.Reporter, reportService *report.ReportService) {
	for {
		select {
		case <-ctx.Done():
			return
		case evt := <-reportService.EventChannel:
			runnerReporter.EventChannel <- evt
		}
	}
}

func handleReportLog(policy Policy, pluginLog io.Writer, runnerReporter *reporter.Reporter) {
	riskLevelFilter := make(map[string]struct{})
	for _, level := range policy.RiskLevelFilter {
		riskLevelFilter[level] = struct{}{}
	}

	enable := false
	events, _ := runnerReporter.GetEvents()
	for _, event := range events {
		if _, ok := riskLevelFilter[toLevelStr(event.Level)]; !ok {
			continue
		}
		enable = true
	}

	if enable {
		if err := runnerReporter.Write(pluginLog); err != nil {
			log.Warn(err)
		}
	}
}

func handleReportAlert(policy Policy, runnerReporter *reporter.Reporter) {
	riskLevelFilter := make(map[string]struct{})
	for _, level := range policy.RiskLevelFilter {
		riskLevelFilter[level] = struct{}{}
	}
	events, _ := runnerReporter.GetEvents()
	for _, event := range events {
		if _, ok := riskLevelFilter[toLevelStr(event.Level)]; !ok {
			continue
		}

		if policy.Alert {
			log.Warn(fmt.Sprintf("Action %s has risks!", policy.Action))
		}
	}
}
