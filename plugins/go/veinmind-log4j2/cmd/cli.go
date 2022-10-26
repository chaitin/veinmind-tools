package main

import (
	"encoding/json"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-log4j2/pkg/scanner"
	"os"
	"time"
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

func scanImage(c *cmd.Command, image api.Image) error {
	var result []*scanner.Result
	defer func() {
		err := image.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	err := scanner.ScanImage(image, &result)
	if err != nil {
		log.Error("Scan Image Error")
		return err
	}

	detail, err := json.Marshal(result)
	if err != nil {
		log.Error("Marshal Results Error")
		return err
	}

	if len(result) > 0 {
		reportEvent := report.ReportEvent{
			ID:         image.ID(),
			Time:       time.Now(),
			Level:      report.Critical,
			DetectType: report.Image,
			EventType:  report.Risk,
			GeneralDetails: []report.GeneralDetail{
				detail,
			},
		}

		err = report.DefaultReportClient(report.WithDisableLog()).Report(reportEvent)
		return err
	}

	return nil
}

func scanContainer(c *cmd.Command, container api.Container) error {
	var result []*scanner.Result
	defer func() {
		err := container.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	err := scanner.ScanContainer(container, &result)
	if err != nil {
		log.Error("Scan Image Error")
		return err
	}

	detail, err := json.Marshal(result)
	if err != nil {
		log.Error("Marshal Results Error")
		return err
	}

	if len(result) > 0 {
		reportEvent := report.ReportEvent{
			ID:         container.ID(),
			Time:       time.Now(),
			Level:      report.Critical,
			DetectType: report.Container,
			EventType:  report.Risk,
			GeneralDetails: []report.GeneralDetail{
				detail,
			},
		}

		err = report.DefaultReportClient(report.WithDisableLog()).Report(reportEvent)
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(cmd.MapImageCommand(scanImageCmd, scanImage))
	rootCmd.AddCommand(cmd.MapContainerCommand(scanContainerCmd, scanContainer))
	rootCmd.AddCommand(cmd.NewInfoCommand(plugin.Manifest{
		Name:        "veinmind-log4j2",
		Author:      "veinmind-team",
		Description: "veinmind-log4j2 scanner image which has log4j jar vulnerable with CVE-2021-44228",
	}))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
