package authz

import (
	"context"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/reporter"
)

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
