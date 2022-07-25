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

func handleReportEvents(eventListCh <-chan []reporter.ReportEvent, policy Policy,
	pluginLog io.Writer) {
	riskLevelFilter := make(map[string]struct{})
	for _, level := range policy.RiskLevelFilter {
		riskLevelFilter[level] = struct{}{}
	}
	select {
	case events := <-eventListCh:
		filter := true
		for _, event := range events {
			if _, ok := riskLevelFilter[toLevelStr(event.Level)]; !ok {
				continue
			}

			filter = false
		}

		if !filter {
			if policy.Alert {
				log.Warn(fmt.Sprintf("Action %s has risks!", policy.Action))
			}
		}

		if err := reporter.WriteEvents2Log(events, pluginLog); err != nil {
			log.Warn(err)
		}
	}
}
