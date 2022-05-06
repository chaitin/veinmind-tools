package main

import (
	"os"
	"time"

	"github.com/aquasecurity/fanal/types"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	scanner "github.com/chaitin/veinmind-tools/veinmind-asset/analyzer"
	"github.com/chaitin/veinmind-tools/veinmind-asset/model"
	"github.com/chaitin/veinmind-tools/veinmind-asset/utils"
	"github.com/chaitin/veinmind-tools/veinmind-common/go/service/report"
	"github.com/spf13/cobra"
)

var results = []model.ScanImageResult{}
var scanStart = time.Now()
var rootCmd = &cmd.Command{}

var scanCmd = &cmd.Command{
	Use:   "scan",
	Short: "Scan image asset",
	PostRun: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")
		verbose, _ := cmd.Flags().GetBool("verbose")
		pkgType, _ := cmd.Flags().GetString("type")
		spend := time.Since(scanStart)

		utils.OutputStream(spend, results, func() error {
			switch format {
			case utils.STDOUT:
				if err := utils.OutputStdout(verbose, pkgType, results); err != nil {
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
		})
	},
}

func scan(c *cmd.Command, image api.Image) error {
	threads, _ := c.Flags().GetInt64("threads")
	res, err := scanner.ScanImage(image, threads)
	if err != nil {
		log.Error("Scan Image Error")
		return err
	}

	results = append(results, res)

	// result event
	if res.PackageTotal > 0 || res.ApplicationTotal > 0 {
		details := []report.AlertDetail{}
		for _, pkg := range res.Packages {
			details = append(details, report.AlertDetail{
				AssetDetail: &report.AssetDetail{
					Type:            "os-pkg",
					Name:            pkg.Name,
					Version:         pkg.Version,
					Release:         pkg.Release,
					Epoch:           pkg.Epoch,
					Arch:            pkg.Arch,
					SrcName:         pkg.SrcName,
					SrcVersion:      pkg.SrcVersion,
					SrcRelease:      pkg.SrcRelease,
					SrcEpoch:        pkg.SrcEpoch,
					Modularitylabel: pkg.Modularitylabel,
					Indirect:        pkg.Indirect,
					License:         pkg.License,
					Layer: func() string {
						if pkg.Layer != (types.Layer{}) {
							return pkg.Layer.DiffID
						} else {
							return ""
						}
					}(),
					FilePath: pkg.FilePath,
				},
			})
		}
		for _, info := range res.Applications {
			for _, lib := range info.Libraries {
				details = append(details, report.AlertDetail{
					AssetDetail: &report.AssetDetail{
						Type:            info.Type,
						Name:            lib.Name,
						Version:         lib.Version,
						Release:         lib.Release,
						Epoch:           lib.Epoch,
						Arch:            lib.Arch,
						SrcName:         lib.SrcName,
						SrcVersion:      lib.SrcVersion,
						SrcRelease:      lib.SrcRelease,
						SrcEpoch:        lib.SrcEpoch,
						Modularitylabel: lib.Modularitylabel,
						Indirect:        lib.Indirect,
						License:         lib.License,
						Layer: func() string {
							if lib.Layer != (types.Layer{}) {
								return lib.Layer.DiffID
							} else {
								return ""
							}
						}(),
						FilePath: func() string {
							if lib.FilePath != "" {
								return lib.FilePath
							} else {
								return info.FilePath
							}
						}(),
					},
				})
			}
		}
		reportEvent := report.ReportEvent{
			ID:           image.ID(),
			Time:         time.Now(),
			Level:        report.Low,
			DetectType:   report.Image,
			EventType:    report.Info,
			AlertType:    report.Asset,
			AlertDetails: details,
		}
		err = report.DefaultReportClient().Report(reportEvent)
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(cmd.MapImageCommand(scanCmd, scan))
	rootCmd.AddCommand(cmd.NewInfoCommand(plugin.Manifest{
		Name:        "veinmind-asset",
		Author:      "veinmind-team",
		Description: "veinmind-asset scanner image os/pkg/app/ info",
	}))
	scanCmd.Flags().Int64P("threads", "t", 10, "scan file threads")
	scanCmd.Flags().StringP("format", "f", "stdout", "epxort file format")
	scanCmd.Flags().BoolP("verbose", "v", false, "show detail Info")
	scanCmd.Flags().String("type", "all", "show specify type detail Info")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
