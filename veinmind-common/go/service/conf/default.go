package conf

import (
	"context"
	"errors"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/libveinmind/go/plugin/service"
	"golang.org/x/sync/errgroup"
	"sync"
)

var (
	defaultOnce   sync.Once
	defaultError  error
	defaultClient *confClient
)

func DefaultConfClient() *confClient {
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
			var pull func(ns PluginConfNS) ([]byte, error)
			service.GetService(Namespace, "pull", &pull)
			group, ctx := errgroup.WithContext(context.Background())

			defaultClient = &confClient{
				ctx:    ctx,
				group:  group,
				Pull:   pull,
			}
		} else {
			group, ctx := errgroup.WithContext(context.Background())

			defaultClient = &confClient{
				ctx:   ctx,
				group: group,
				Pull: func(ns PluginConfNS) ([]byte, error) {
					return nil, errors.New("conf: please use pull action in service mode")
				},
			}
		}
	})

	if defaultError != nil {
		log.Error(defaultError)
	}
	return defaultClient
}