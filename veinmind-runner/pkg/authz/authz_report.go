package authz

import (
	"fmt"
	"io"

	gomail "gopkg.in/mail.v2"

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

func handleHarborWebhookReportEvents(eventListCh <-chan []reporter.ReportEvent, hpolicy HarborPolicy,
	pluginLog io.Writer, mailconf MailConf) {
	filter, events := processReportEvents(eventListCh, hpolicy.Policy, pluginLog)
	if filter {
		if hpolicy.Alert {
			log.Warn(fmt.Sprintf("Action %s has risks!", hpolicy.Action))
		}
		if hpolicy.SendMail {
			err := sendReport2Mail(events, mailconf)
			if err != nil {
				log.Error(err)
			}
		}
	}
	if err := reporter.WriteEvents2Log(events, pluginLog); err != nil {
		log.Warn(err)
	}
}
func sendReport2Mail(events []reporter.ReportEvent, mailconf MailConf) error {
	d := gomail.NewDialer(mailconf.Host, mailconf.Port, mailconf.Name, mailconf.Password)
	m := gomail.NewMessage()
	m.SetHeader("From", mailconf.Name)
	m.SetHeader("To", mailconf.Subscriber...)
	m.SetHeader("Subject", "Harbor webhook Report")
	m.SetBody("text/plain", fmt.Sprintf("%#v", events))
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
