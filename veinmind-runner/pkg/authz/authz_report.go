package authz

import (
	"context"
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
