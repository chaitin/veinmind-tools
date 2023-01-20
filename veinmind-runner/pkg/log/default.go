package log

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	defaultModule *Module
	defaultOnce   sync.Once
)

func DefaultModule() *Module {
	defaultOnce.Do(func() {
		logger := logrus.New()

		// set formatter
		logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			ForceQuote:      true,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		})

		// set level
		logger.SetLevel(logrus.InfoLevel)

		// entry
		entry := logger.WithFields(map[string]interface{}{
			"module": "default",
		})

		defaultModule = &Module{
			name:  "default",
			Entry: entry,
		}
	})

	return defaultModule
}
