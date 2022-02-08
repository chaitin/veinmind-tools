package common_cli

import (
	"github.com/chaitin/veinmind-tools/veinmind-malicious/database/model"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/embed"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/scanner"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/scanner/scanner_common"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/sdk/common"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/sdk/common/report"
	"github.com/urfave/cli/v2"
	"strings"
	"time"
)

var App = &cli.App{
	Name:  "Veinmind-Imagescan",
	Usage: "Veinmind-Imagescan is a image scanner",
	Commands: []*cli.Command{
		{
			Name:  "scan",
			Usage: "scan image, image e.g. ubuntu:latest",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "name",
					Value:   "report",
					Usage:   "report name for scan report",
					Aliases: []string{"n"},
				},
				&cli.StringFlag{
					Name:    "format",
					Value:   "html",
					Usage:   "report format for scan report",
					Aliases: []string{"f"},
				},
				&cli.StringFlag{
					Name:    "engine",
					Value:   "dockerd",
					Usage:   "scan engine e.g. dockerd",
					Aliases: []string{"e"},
				},
				&cli.StringSliceFlag{
					Name:    "plugins",
					Value:   cli.NewStringSlice(scanner.ScanPluginsName...),
					Usage:   "scan plugin",
					Aliases: []string{"p"},
				},
			},
			Action: func(c *cli.Context) error {
				var scanReport model.ReportData
				var err error

				// 记录扫描开始时间
				scanStart := time.Now()

				// 引擎类型
				var engineType scanner_common.ScanEngineType
				switch c.String("engine") {
				case "dockerd":
					engineType = scanner_common.Dockerd
				case "containerd":
					engineType = scanner_common.Containerd
				}

				if c.Args().First() != "" {
					image := c.Args().First()
					scanReport, err = scanner.Scan(scanner_common.ScanOption{
						EngineType:    engineType,
						ImageName:     image,
						EnablePlugins: c.StringSlice("plugins"),
					})
					if err != nil {
						common.Log.Error(err)
						return err
					}
				} else {
					scanReport, err = scanner.Scan(scanner_common.ScanOption{
						EngineType:    engineType,
						EnablePlugins: c.StringSlice("plugins"),
					})
					if err != nil {
						common.Log.Error(err)
						return err
					}
				}

				// 计算扫描数据
				spend := time.Since(scanStart)
				scanReport.ScanSpendTime = spend.String()
				scanReport.ScanStartTime = scanStart.Format("2006-01-02 15:04:05")
				report.CalculateScanReportCount(&scanReport)
				report.SortScanReport(&scanReport)

				// 输出扫描结果
				format := c.String("format")
				name := c.String("name")
				name = strings.Join([]string{name, format}, ".")

				switch format {
				case report.HTML:
					report.OutputHTML(scanReport, name)
				case report.JSON:
					report.OutputJSON(scanReport, name)
				case report.CSV:
					report.OutputCSV(scanReport, name)
				}

				return nil
			},
		},
		{
			Name:  "extract",
			Usage: "extract config file to disk",
			Action: func(c *cli.Context) error {
				embed.ExtractAll()
				return nil
			},
		},
	},
}
