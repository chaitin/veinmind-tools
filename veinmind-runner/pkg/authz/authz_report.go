package authz

import (
	"fmt"
	"io"

	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/reporter"
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

func processReportEvents(eventListCh <-chan []reporter.ReportEvent, policy Policy,
	pluginLog io.Writer) (reportFlag bool, results []reporter.ReportEvent) {
	riskLevelFilter := make(map[string]struct{})
	for _, level := range policy.RiskLevelFilter {
		riskLevelFilter[level] = struct{}{}
	}
	select {
	case events := <-eventListCh:
		for _, event := range events {
			if _, ok := riskLevelFilter[toLevelStr(event.Level)]; !ok {
				continue
			}
			reportFlag = true
			results = append(results, event)
		}
	}
	return reportFlag, results
}
func handleDockerPluginReportEvents(eventListCh <-chan []reporter.ReportEvent, bpolicy Policy,
	pluginLog io.Writer) {
	filter, events := processReportEvents(eventListCh, bpolicy, pluginLog)
	if filter {
		if bpolicy.Alert {
			log.Warn(fmt.Sprintf("Action %s has risks!", bpolicy.Action))
		}
	}
	if err := reporter.WriteEvents2Log(events, pluginLog); err != nil {
		log.Warn(err)
	}
}
