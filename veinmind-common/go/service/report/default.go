package report

import (
	"context"
	"encoding/json"
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

// PluginOption use for plugin standalone version (without host)
type PluginOption func(r *reportClient) (*reportClient, error)

func WithDisableLog() PluginOption {
	return func(r *reportClient) (*reportClient, error) {
		r.Report = func(event ReportEvent) error {
			return nil
		}

		return r, nil
	}
}

func DefaultReportClient(pOpts ...PluginOption) *reportClient {
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
			var report func(ReportEvent) (error)
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
				Report: func(evt ReportEvent) error {
					evtBytes, err := json.MarshalIndent(evt, "", "	")
					if err != nil {
						return err
					}
					log.Warn(string(evtBytes))
					return nil
				},
			}

			for _, opt := range pOpts {
				d, err := opt(defaultClient)
				if err != nil {
					log.Error(err)
					continue
				}

				defaultClient = d
			}
		}
	})

	if defaultError != nil {
		panic(defaultError)
	}
	return defaultClient
}
