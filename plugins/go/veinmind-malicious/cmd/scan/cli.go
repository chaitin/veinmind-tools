package main

import (
	_ "net/http/pprof"
	"os"
	"sync"
	"syscall"
	"time"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-common-go/service/report/event"

	_ "github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/config"
	_ "github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/database"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/database/model"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/embed"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/scanner/malicious"
)

var (
	reportData    = model.ReportData{}
	reportLock    sync.Mutex
	ReportService = &report.Service{}
	rootCmd       = &cmd.Command{}
	scanCmd       = &cmd.Command{Use: "scan"}
	extractCmd    = &cmd.Command{
		Use:   "extract",
		Short: "Extract config file",
		RunE: func(cmd *cmd.Command, args []string) error {
			embed.ExtractAll()
			return nil
		},
	}
	scanImageCmd = &cmd.Command{
		Use:   "image",
		Short: "Scan image malicious files",
	}
)

func scan(_ *cmd.Command, image api.Image) error {

	result, err := malicious.Scan(image)
	if err != nil {
		log.Error(err)
		return nil
	}
	reportLock.Lock()
	reportData.ScanImageResult = append(reportData.ScanImageResult, result)
	reportLock.Unlock()

	// result event
	if result.MaliciousFileCount > 0 {
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
					reportEvent := event.Event{
						&event.BasicInfo{
							ID:         image.ID(),
							Object:     event.NewObject(image),
							Time:       time.Now(),
							Level:      event.High,
							DetectType: event.Image,
							EventType:  event.Risk,
							AlertType:  event.MaliciousFile,
						},
						event.NewDetailInfo(&event.MaliciousFileDetail{
							Engine:        mr.Engine,
							MaliciousName: mr.Description,
							FileDetail: event.FileDetail{
								Path: mr.RelativePath,
								Perm: fStat.Mode(),
								Size: fStat.Size(),
								Gid:  int64(fSys.Gid),
								Uid:  int64(fSys.Uid),
								Ctim: fSys.Ctim.Sec,
								Mtim: fSys.Mtim.Sec,
								Atim: fSys.Atim.Sec,
							},
						}),
					}
					err = ReportService.Client.Report(&reportEvent)
					if err != nil {
						log.Error(err)
						continue
					}
				}
			}
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.AddCommand(report.MapReportCmd(cmd.MapImageCommand(scanImageCmd, scan), ReportService))
	rootCmd.AddCommand(extractCmd)
	rootCmd.AddCommand(cmd.NewInfoCommand(plugin.Manifest{
		Name:        "veinmind-malicious",
		Author:      "veinmind-team",
		Description: "veinmind-malicious scanner image malicious file",
	}))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
