package main

import (
	"os"
	"time"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-common-go/service/report/event"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-trace/pkg/analyzer"
)

var reportService = &report.Service{}
var rootCmd = &cmd.Command{}

var scanCmd = &cmd.Command{
	Use:   "scan",
	Short: "scan command",
}

var scanContainerCmd = &cmd.Command{
	Use:   "container",
	Short: "container command",
}

// scanContainer is func that used to do some action with container
// you can write your plugin scan code here
func scanContainer(c *cmd.Command, container api.Container) error {
	defer func() {
		err := container.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	result := make([]*event.TraceEvent, 0)
	for _, a := range analyzer.Group {
		a.Scan(container)
		result = append(result, a.Result()...)
	}

	for _, e := range result {
		reportEvent := &event.Event{
			BasicInfo: &event.BasicInfo{
				ID:         container.ID(), // container id info
				Object:     event.NewObject(container),
				Time:       time.Now(),      // report time, usually use time.Now
				Level:      e.Level,         // report event level
				DetectType: event.Container, // report scan object type
				AlertType:  event.TraceRisk, // report alert type, we provide some clearly types of security events,
				EventType:  event.Risk,      // report event type: Risk/Invasion/Info
			},
			DetailInfo: &event.DetailInfo{
				AlertDetail: e,
			},
		}
		err := reportService.Client.Report(reportEvent)
		if err != nil {
			log.Error(err)
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.AddCommand(report.MapReportCmd(cmd.MapContainerCommand(scanContainerCmd, scanContainer), reportService))
	rootCmd.AddCommand(cmd.NewInfoCommand(plugin.Manifest{
		Name:        "veinmind-trace",
		Author:      "veinmind-team",
		Description: "veinmind-trace Description",
	}))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
