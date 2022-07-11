package main

import (
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report"
	scanner "github.com/chaitin/veinmind-tools/plugins/go/veinmind-asset/analyzer"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-asset/model"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-asset/utils"
	"github.com/spf13/cobra"
	"os"
	"time"
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
	defer func() {
		err := image.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	threads, _ := c.Flags().GetInt64("threads")
	res, err := scanner.ScanImage(image, threads)
	if err != nil {
		log.Error("Scan Image Error")
		return err
	}

	results = append(results, res)

	assetDetail := &report.AssetDetail{
		OS: report.AssetOSDetail{
			Family: res.ImageOSInfo.Family,
			Name:   res.ImageOSInfo.Name,
			Eosl:   res.ImageOSInfo.Eosl,
		},
		PackageInfos: func() []report.AssetPackageDetails {
			assetPackageDetailsList := []report.AssetPackageDetails{}
			assetPackageDetails := []report.AssetPackageDetail{}

			for _, pkgInfo := range res.PackageInfos {
				for _, pkg := range pkgInfo.Packages {
					assetPackageDetails = append(assetPackageDetails, report.AssetPackageDetail{
						Name:       pkg.Name,
						Version:    pkg.Version,
						Release:    pkg.Release,
						Epoch:      pkg.Epoch,
						Arch:       pkg.Arch,
						SrcName:    pkg.SrcName,
						SrcEpoch:   pkg.SrcEpoch,
						SrcRelease: pkg.SrcRelease,
						SrcVersion: pkg.SrcVersion,
					})
				}
				assetPackageDetailsList = append(assetPackageDetailsList, report.AssetPackageDetails{
					FilePath: pkgInfo.FilePath,
					Packages: assetPackageDetails,
				})
				assetPackageDetails = []report.AssetPackageDetail{}
			}

			return assetPackageDetailsList
		}(),
		Applications: func() []report.AssetApplicationDetails {
			assetApplicationDetailsList := []report.AssetApplicationDetails{}
			assetPackageDetails := []report.AssetPackageDetail{}

			for _, app := range res.Applications {
				for _, pkg := range app.Libraries {
					assetPackageDetails = append(assetPackageDetails, report.AssetPackageDetail{
						Name:       pkg.Name,
						Version:    pkg.Version,
						Release:    pkg.Release,
						Epoch:      pkg.Epoch,
						Arch:       pkg.Arch,
						SrcName:    pkg.SrcName,
						SrcEpoch:   pkg.SrcEpoch,
						SrcRelease: pkg.SrcRelease,
						SrcVersion: pkg.SrcVersion,
					})
				}
				assetApplicationDetailsList = append(assetApplicationDetailsList, report.AssetApplicationDetails{
					Type:     app.Type,
					FilePath: app.FilePath,
					Packages: assetPackageDetails,
				})
				assetPackageDetails = []report.AssetPackageDetail{}
			}

			return assetApplicationDetailsList
		}(),
	}

	reportEvent := report.ReportEvent{
		ID:         image.ID(),
		Time:       time.Now(),
		Level:      report.None,
		DetectType: report.Image,
		EventType:  report.Info,
		AlertType:  report.Asset,
		AlertDetails: []report.AlertDetail{
			{
				AssetDetail: assetDetail,
			},
		},
	}
	err = report.DefaultReportClient(report.WithDisableLog()).Report(reportEvent)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(cmd.MapImageCommand(scanCmd, scan))
	rootCmd.AddCommand(cmd.NewInfoCommand(plugin.Manifest{
		Name:        "veinmind-asset",
		Author:      "veinmind-team",
		Description: "veinmind-asset scanner image os/pkg/app info",
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
