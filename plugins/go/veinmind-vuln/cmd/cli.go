package main

import (
	"os"
	"time"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report"
	scanner "github.com/chaitin/veinmind-tools/plugins/go/veinmind-vuln/analyzer"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-vuln/model"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-vuln/utils"
)

var (
	results   = []model.ScanResult{}
	scanStart = time.Now()
	rootCmd   = &cmd.Command{}
	post      = func(cmd *cmd.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")
		verbose, _ := cmd.Flags().GetBool("verbose")
		onlyAsset, _ := cmd.Flags().GetBool("only-asset")
		pkgType, _ := cmd.Flags().GetString("type")
		spend := time.Since(scanStart)

		utils.OutputStream(spend, results, func() error {
			switch format {
			case utils.STDOUT:
				if err := utils.OutputStdout(verbose, onlyAsset, pkgType, results); err != nil {
					log.Error("Stdout error", err)
					return err
				}
			case utils.JSON:
				if err := utils.OutputJSON(results); err != nil {
					log.Error("Export Results JSON False", err)
					return err
				}
			case utils.CSV:
				if err := utils.OutputCSV(results); err != nil {
					log.Error("Export Results CSV False", err)
					return err
				}
			}
			return nil
		}, onlyAsset)
	}
	pre = func(cmd *cmd.Command, args []string) {
		if _, err := os.Open("./data"); os.IsNotExist(err) {
			_ = os.Mkdir("./data", 0600)
		}
	}

	scanCmd = &cmd.Command{
		Use:   "scan",
		Short: "Scan asset and vulns",
	}

	scanImageCmd = &cmd.Command{
		Use:     "image",
		Short:   "Scan image asset/vulns",
		PreRun:  pre,
		PostRun: post,
	}

	scanContainerCmd = &cmd.Command{
		Use:     "container",
		Short:   "Scan container asset/vulns",
		PreRun:  pre,
		PostRun: post,
	}
)

func scanImage(c *cmd.Command, image api.Image) error {
	defer func() {
		err := image.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	threads, _ := c.Flags().GetInt64("threads")
	onlyAsset, _ := c.Flags().GetBool("only-asset")
	verbose, _ := c.Flags().GetBool("verbose")
	res, err := scanner.ScanImage(image, threads)
	if err != nil {
		log.Error("Scan Image Error")
		return err
	}
	// cve
	if !onlyAsset {
		scanner.ScanOSV(&res, verbose)
	}
	// first format, then add
	results = append(results, res)

	if onlyAsset {
		reportEvent := report.ReportEvent{
			ID:         image.ID(),
			Time:       time.Now(),
			Level:      report.None,
			DetectType: report.Image,
			EventType:  report.Info,
			AlertType:  report.Asset,
			AlertDetails: []report.AlertDetail{
				{
					AssetDetail: scanner.TransferAsset(res),
				},
			},
		}
		err = report.DefaultReportClient(report.WithDisableLog()).Report(reportEvent)
		if err != nil {
			return err
		}
	}
	if res.CveTotal > 0 {
		reportEvent := report.ReportEvent{
			ID:         image.ID(),
			Time:       time.Now(),
			Level:      report.High,
			DetectType: report.Image,
			EventType:  report.Risk,
			AlertType:  report.Vulnerability,
			GeneralDetails: []report.GeneralDetail{
				scanner.TransferVuln(res),
			},
		}
		err = report.DefaultReportClient(report.WithDisableLog()).Report(reportEvent)
		if err != nil {
			return err
		}
	}

	return nil
}

func scanContainer(c *cmd.Command, container api.Container) error {
	defer func() {
		err := container.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	threads, _ := c.Flags().GetInt64("threads")
	onlyAsset, _ := c.Flags().GetBool("only-asset")
	verbose, _ := c.Flags().GetBool("verbose")
	res, err := scanner.ScanContainer(container, threads)
	if err != nil {
		log.Error("Scan Image Error")
		return err
	}

	if !onlyAsset {
		scanner.ScanOSV(&res, verbose)
	}
	// first format, then add
	results = append(results, res)

	if onlyAsset {
		reportEvent := report.ReportEvent{
			ID:         container.ID(),
			Time:       time.Now(),
			Level:      report.None,
			DetectType: report.Container,
			EventType:  report.Info,
			AlertType:  report.Asset,
			AlertDetails: []report.AlertDetail{
				{
					AssetDetail: scanner.TransferAsset(res),
				},
			},
		}
		err = report.DefaultReportClient(report.WithDisableLog()).Report(reportEvent)
		if err != nil {
			return err
		}
	}

	if res.CveTotal > 0 {
		reportEvent := report.ReportEvent{
			ID:         container.ID(),
			Time:       time.Now(),
			Level:      report.High,
			DetectType: report.Container,
			EventType:  report.Risk,
			AlertType:  report.Vulnerability,
			GeneralDetails: []report.GeneralDetail{
				scanner.TransferVuln(res),
			},
		}
		err = report.DefaultReportClient(report.WithDisableLog()).Report(reportEvent)
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.AddCommand(cmd.MapImageCommand(scanImageCmd, scanImage))
	scanCmd.AddCommand(cmd.MapContainerCommand(scanContainerCmd, scanContainer))

	rootCmd.AddCommand(cmd.NewInfoCommand(plugin.Manifest{
		Name:        "veinmind-vuln",
		Author:      "veinmind-team",
		Description: "veinmind-vuln scanner image os/pkg/app info and vulns",
	}))
	scanCmd.PersistentFlags().Int64P("threads", "t", 5, "scan file threads")
	scanCmd.PersistentFlags().StringP("format", "f", "stdout", "epxort file format")
	scanCmd.PersistentFlags().BoolP("verbose", "v", false, "show detail Info")
	scanCmd.PersistentFlags().String("type", "all", "show specify type detail Info")
	scanCmd.PersistentFlags().Bool("only-asset", false, "only scan asset info")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
