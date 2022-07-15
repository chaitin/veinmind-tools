package main

import (
	"github.com/biter777/processex"
	"github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	_ "github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/config"
	_ "github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/database"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/database/model"
	_ "github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/database/model"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/embed"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/scanner/malicious"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/sdk/common/report"
	reportService "github.com/chaitin/veinmind-tools/veinmind-common/go/service/report"
	"github.com/spf13/cobra"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"syscall"
	"time"
)

var reportData = model.ReportData{}
var reportLock sync.Mutex
var scanStart = time.Now()
var clamdPid = -1

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
		manual, _ := cmd.Flags().GetBool("manual")
		// the permission of use run the clamAV
		if !manual {
			log.Info("using local clamAV")
			// make sure local is not running a clamAV
			_, _, err := processex.FindByName("clamd")
			if err == processex.ErrNotFound {
				log.Info("start local clamAV")
				clam := exec.Command("clamd")
				err := clam.Run()
				//get the clamd PID
				clamdPid = clam.Process.Pid + 1
				if err != nil {
					log.Error("clamAV failed to start: " + err.Error())
				}
			} else if err == nil {
				log.Info("the local clamAV is running")
			} else {
				log.Error("find clamAV ERROR")
			}
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		// 计算扫描数据
		spend := time.Since(scanStart)
		reportData.ScanSpendTime = spend.String()
		reportData.ScanStartTime = scanStart.Format("2006-01-02 15:04:05")
		report.CalculateScanReportCount(&reportData)
		report.SortScanReport(&reportData)

		format, _ := cmd.Flags().GetString("format")
		name, _ := cmd.Flags().GetString("name")
		outputPath, _ := cmd.Flags().GetString("output")
		name = strings.Join([]string{name, format}, ".")
		fpath := path.Join(outputPath, name)

		switch format {
		case report.HTML:
			report.OutputHTML(reportData, fpath)
		case report.JSON:
			report.OutputJSON(reportData, fpath)
		case report.CSV:
			report.OutputCSV(reportData, fpath)
		}

		if clamdPid != -1 {
			log.Info("close clamAV server")
			err := syscall.Kill(-clamdPid, syscall.SIGKILL)
			if err != nil {
				log.Error(err)
			}
		}
	},
}

func scan(c *cmd.Command, image api.Image) error {
	// clamAV host and port
	clamdHost, _ := c.Flags().GetString("remote")
	clamdPort, _ := c.Flags().GetString("port")

	result, err := malicious.Scan(image, clamdHost, clamdPort)
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
	scanCmd.Flags().BoolP("manual", "m", false, "whether need to manually start clamAV")
	scanCmd.Flags().StringP("remote", "r", "127.0.0.1", "host of ClamAV")
	scanCmd.Flags().StringP("port", "p", "3310", "port of ClamAV")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
