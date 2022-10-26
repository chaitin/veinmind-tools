package main

import (
	"os"
	"time"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report"
)

var rootCmd = &cmd.Command{}
var scanImageCmd = &cmd.Command{
	Use:   "scan-image",
	Short: "scan image command",
}

var scanContainerCmd = &cmd.Command{
	Use:   "scan-container",
	Short: "scan container command",
}

// scanImage is func that used to do some action with Images
// you can write your plugin scan code here
func scanImage(c *cmd.Command, image api.Image) error {
	defer func() {
		err := image.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	// do something here
	log.Info(image.ID())

	// if you want display at runner report, you should send your result to report event
	reportEvent := report.ReportEvent{
		ID:             image.ID(),               // image id info
		Time:           time.Now(),               // report time, usually use time.Now
		Level:          report.None,              // report event level
		DetectType:     report.Image,             // report scan object type
		EventType:      report.Info,              // report event type: Risk/Invasion/Info
		AlertType:      report.Asset,             // report alert type, we provide some clearly types of security events,
		AlertDetails:   []report.AlertDetail{},   // add report detail data in there
		GeneralDetails: []report.GeneralDetail{}, // if your report event does not in alert type, you can use GeneralDetails type which consists of json bytes
	}
	err := report.DefaultReportClient(report.WithDisableLog()).Report(reportEvent)
	if err != nil {
		return err
	}

	return nil
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
	// do something here
	log.Info(container.ID())

	// if you want display at runner report, you should send your result to report event
	reportEvent := report.ReportEvent{
		ID:             container.ID(),           // image id info
		Time:           time.Now(),               // report time, usually use time.Now
		Level:          report.None,              // report event level
		DetectType:     report.Container,         // report scan object type
		EventType:      report.Info,              // report event type: Risk/Invasion/Info
		AlertType:      report.Asset,             // report alert type, we provide some clearly types of security events,
		AlertDetails:   []report.AlertDetail{},   // add report detail data in there
		GeneralDetails: []report.GeneralDetail{}, // if your report event does not in alert type, you can use GeneralDetails type which consists of json bytes
	}
	err := report.DefaultReportClient(report.WithDisableLog()).Report(reportEvent)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(cmd.MapImageCommand(scanImageCmd, scanImage))
	rootCmd.AddCommand(cmd.MapContainerCommand(scanContainerCmd, scanContainer))
	rootCmd.AddCommand(cmd.NewInfoCommand(plugin.Manifest{
		Name:        "veinmind-example",
		Author:      "veinmind-team",
		Description: "veinmind-example Description",
	}))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
