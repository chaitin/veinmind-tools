package authz

import (
	"encoding/json"
	"fmt"
	"github.com/chaitin/veinmind-common-go/service/report/event"
	"io"

	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/log"
)

func toLevelStr(level event.Level) string {
	switch level {
	case event.Low:
		return "Low"
	case event.Medium:
		return "Medium"
	case event.High:
		return "High"
	case event.Critical:
		return "Critical"
	}

	return "None"
}

func processReportEvents(eventListCh <-chan []*event.Event, policy Policy,
	pluginLog io.Writer) (reportFlag bool, results []*event.Event) {
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

func handleDockerPluginReportEvents(eventListCh <-chan []*event.Event, bpolicy Policy,
	pluginLog io.Writer) {
	filter, events := processReportEvents(eventListCh, bpolicy, pluginLog)
	if filter {
		if bpolicy.Alert {
			log.GetModule(log.AuthzModuleKey).Warn(fmt.Sprintf("action %s has risks!", bpolicy.Action))
		}
	}
	if err := WriteEvents2Log(events, pluginLog); err != nil {
		log.GetModule(log.AuthzModuleKey).Warn(err)
	}
}

func WriteEvents2Log(events []*event.Event, writer io.Writer) error {
	if len(events) == 0 {
		return nil
	}

	eventsBytes, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return err
	}

	_, err = writer.Write(eventsBytes)
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte("\n"))
	return err
}
