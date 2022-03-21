package report

import (
	"context"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/libveinmind/go/plugin/service"
	"golang.org/x/sync/errgroup"
	"sync"
)

var (
	defaultOnce       sync.Once
	defaultError      error
	defaultClient     *reportClient
)

func DefaultReportClient() *reportClient {
	defaultOnce.Do(func() {
		hasService := false
		if service.Hosted() {
			ok, err := service.HasNamespace(Namespace)
			if err != nil {
				defaultError = err
			}
			hasService = ok
		}

		if hasService {
			var report func([]ReportEvent) (error)
			service.GetService(Namespace, "report", &report)
			group, ctx := errgroup.WithContext(context.Background())

			defaultClient = &reportClient{
				ctx: ctx,
				group: group,
				Report: report,
			}
		} else {
			group, ctx := errgroup.WithContext(context.Background())

			defaultClient = &reportClient{
				ctx: ctx,
				group: group,
				Report: func(evts []ReportEvent) error {
					for _, evt := range evts {
						log.Info(evt)
					}

					return nil
				},
			}
		}
	})

	if defaultError != nil {
		panic(defaultError)
	}
	return defaultClient
}
