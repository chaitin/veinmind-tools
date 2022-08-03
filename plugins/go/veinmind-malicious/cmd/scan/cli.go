package main

import (
	"context"
	"github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	_ "github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/config"
	_ "github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/database"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/database/model"
	_ "github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/database/model"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/embed"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/pkg/avutil"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/scanner/malicious"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/sdk/common/report"
	reportService "github.com/chaitin/veinmind-tools/veinmind-common/go/service/report"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	_ "net/http/pprof"
	"os"
	"path"
	"strings"
	"sync"
	"syscall"
	"time"
)

var reportData = model.ReportData{}
var reportLock sync.Mutex
var scanStart = time.Now()
var ctx = context.Background()
var cancel context.CancelFunc
var clamavAutoStart bool
var clamAV *avutil.ClamAVManger

var rootCmd = &cmd.Command{}
var extractCmd = &cmd.Command{
	Use:   "extract",
	Short: "Extract config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		embed.ExtractAll()
		return nil
	},
}
var scanCmd = &cmd.Command{
	Use:   "scan",
	Short: "Scan image malicious files",
	PreRun: func(cmd *cobra.Command, args []string) {
		clamavConf, err := cmd.Flags().GetString("clamav-conf")
		if err != nil {
			log.Error(err)
		}
		clamavExec, err := cmd.Flags().GetString("clamav-exec")
		if err != nil {
			log.Error(err)
		}
		clamavPort, err := cmd.Flags().GetString("clamav-port")
		if err != nil {
			log.Error(err)
		}
		clamavHost, err := cmd.Flags().GetString("clamav-host")
		if err != nil {
			log.Error(err)
		}

		// the flag of Manual run the clamAV
		if clamavAutoStart {
			log.Info("start clamAV ....")
			ctx, cancel = context.WithCancel(ctx)
			options := make([]avutil.ServiceOption, 0)
			options = append(options, avutil.WithHost(clamavHost), avutil.WithPort(clamavPort))
			options = append(options, avutil.WithExec(clamavExec), avutil.WithConf(clamavConf))
			clamAV = avutil.New(ctx, options...)
			go func() {
				err := clamAV.Run()
				if err != nil {
					log.Error(err)
				}
			}()

			g, _ := errgroup.WithContext(ctx)
			g.Go(clamAV.Ready)

			if err = g.Wait(); err != nil {
				log.Error(err)
				return
			}

			go func() {
				err := clamAV.Daemon()
				if err != nil {
					log.Info(err)
				}
			}()
		}
	},

	PostRun: func(cmd *cobra.Command, args []string) {
		// 计算扫描数据
		spend := time.Since(scanStart)
		reportData.ScanSpendTime = spend.String()
		reportData.ScanStartTime = scanStart.Format("2006-01-02 15:04:05")
		report.CalculateScanReportCount(&reportData)
		report.SortScanReport(&reportData)

		format, err := cmd.Flags().GetString("format")
		if err != nil {
			log.Error(err)
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Error(err)
		}
		outputPath, err := cmd.Flags().GetString("output")
		if err != nil {
			log.Error(err)
		}
		name = strings.Join([]string{name, format}, ".")
		fpath := path.Join(outputPath, name)
		if err != nil {
			log.Error(err)
		}

		switch format {
		case report.HTML:
			report.OutputHTML(reportData, fpath)
		case report.JSON:
			report.OutputJSON(reportData, fpath)
		case report.CSV:
			report.OutputCSV(reportData, fpath)
		}

		if clamAV != nil {
			log.Info("close clamAV ....")
			cancel()
			clamAV.Wait()
		}
	},
}

func scan(c *cmd.Command, image api.Image) error {
	clamavHost, err := c.Flags().GetString("clamav-host")
	if err != nil {
		log.Error(err)
	}
	clamavPort, err := c.Flags().GetString("clamav-port")
	if err != nil {
		log.Error(err)
	}

	antiVirusAgent := malicious.AntiVirusEngine{ClamavHost: clamavHost, ClamavPort: clamavPort}
	result, err := malicious.Scan(image, antiVirusAgent)
	if err != nil {
		log.Error(err)
		return nil
	}

	reportLock.Lock()
	reportData.ScanImageResult = append(reportData.ScanImageResult, result)
	reportLock.Unlock()

	// result event
	if result.MaliciousFileCount > 0 {
		details := []reportService.AlertDetail{}
		for _, l := range result.Layers {
			if len(l.MaliciousFileInfos) > 0 {
				for _, mr := range l.MaliciousFileInfos {
					f, err := image.Open(mr.RelativePath)
					if err != nil {
						log.Error(err)
						continue
					}

					fStat, err := f.Stat()
					if err != nil {
						log.Error(err)
						continue
					}
					fSys := fStat.Sys().(*syscall.Stat_t)

					details = append(details, reportService.AlertDetail{
						MaliciousFileDetail: &reportService.MaliciousFileDetail{
							Engine:        mr.Engine,
							MaliciousName: mr.Description,
							FileDetail: reportService.FileDetail{
								Path: mr.RelativePath,
								Perm: fStat.Mode(),
								Size: fStat.Size(),
								Gid:  int64(fSys.Gid),
								Uid:  int64(fSys.Uid),
								Ctim: fSys.Ctim.Sec,
								Mtim: fSys.Mtim.Sec,
								Atim: fSys.Atim.Sec,
							},
						},
					})
				}
			}
		}
		reportEvent := reportService.ReportEvent{
			ID:           image.ID(),
			Level:        reportService.High,
			DetectType:   reportService.Image,
			EventType:    reportService.Risk,
			AlertType:    reportService.MaliciousFile,
			AlertDetails: details,
		}
		err = reportService.DefaultReportClient().Report(reportEvent)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(cmd.MapImageCommand(scanCmd, scan))
	rootCmd.AddCommand(extractCmd)
	rootCmd.AddCommand(cmd.NewInfoCommand(plugin.Manifest{
		Name:        "veinmind-malicious",
		Author:      "veinmind-team",
		Description: "veinmind-malicious scanner image malicious file",
	}))
	scanCmd.Flags().StringP("format", "f", "html", "report format for scan report")
	scanCmd.Flags().StringP("name", "n", "report", "report name for scan report")
	scanCmd.Flags().StringP("output", "o", ".", "output path for report")
	scanCmd.Flags().BoolVarP(&clamavAutoStart, "clamav-automaticly-start", "", false, "whether need to automatically start clamAV")
	scanCmd.Flags().StringP("clamav-host", "", "127.0.0.1", "host of ClamAV")
	scanCmd.Flags().StringP("clamav-port", "", "3310", "port of ClamAV")
	scanCmd.Flags().StringP("clamav-exec", "", "/usr/sbin/clamd", "execution file path of ClamAV")
	scanCmd.Flags().StringP("clamav-conf", "", "/etc/clamav/clamd.conf", "config file path of ClamAV")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
