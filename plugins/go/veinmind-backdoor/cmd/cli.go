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
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-backdoor/utils"
)

var (
	reportService = &report.Service{}
	pluginInfo    = plugin.Manifest{
		Name:        "veinmind-backdoor",
		Author:      "veinmind-team",
		Description: "detect backdoor risk for image&container",
	}
	rootCmd = &cmd.Command{}
	scanCmd = &cmd.Command{
		Use:   "scan",
		Short: "scan mode",
	}
	scanImageCmd = &cmd.Command{
		Use:   "image",
		Short: "scan image backdoor risk",
	}
	scanContainerCmd = &cmd.Command{
		Use:   "container",
		Short: "scan container backdoor risk",
	}
)

func scanImage(c *cmd.Command, image api.Image) error {
	results := utils.ImagesScanRun(image)
	for _, result := range results {
		ReportEvent := &event.Event{
			BasicInfo: &event.BasicInfo{
				ID:         image.ID(),
				Time:       time.Now(),
				Level:      event.High,
				Source:     pluginInfo.Name,
				Object:     event.NewObject(image),
				EventType:  event.Risk,
				DetectType: event.Image,
				AlertType:  event.Escape,
			},
			DetailInfo: &event.DetailInfo{
				AlertDetail: &event.BackdoorDetail{
					FileDetail: event.FileDetail{
						Path:  result.FileDetail.Path,
						Perm:  result.FileDetail.Perm,
						Size:  result.FileDetail.Size,
						Gname: result.FileDetail.Gname,
						Gid:   result.FileDetail.Gid,
						Uid:   result.FileDetail.Uid,
						Uname: result.FileDetail.Uname,
						Ctim:  result.FileDetail.Ctim,
						Mtim:  result.FileDetail.Mtim,
						Atim:  result.FileDetail.Atim,
					},
					Content:     result.Content,
					Description: result.Description,
				},
			},
		}
		err := reportService.Client.Report(ReportEvent)
		if err != nil {
			log.Error(err)
			continue
		}
	}

	return nil
}

func scanContainer(c *cmd.Command, container api.Container) error {
	results := utils.ContainersScanRun(container)
	for _, result := range results {
		ReportEvent := &event.Event{
			BasicInfo: &event.BasicInfo{
				ID:         container.ID(),
				Time:       time.Now(),
				Source:     pluginInfo.Name,
				Level:      event.High,
				Object:     event.NewObject(container),
				EventType:  event.Risk,
				DetectType: event.Container,
				AlertType:  event.Escape,
			},
			DetailInfo: &event.DetailInfo{
				AlertDetail: &event.BackdoorDetail{
					FileDetail: event.FileDetail{
						Path:  result.FileDetail.Path,
						Perm:  result.FileDetail.Perm,
						Size:  result.FileDetail.Size,
						Gname: result.FileDetail.Gname,
						Gid:   result.FileDetail.Gid,
						Uid:   result.FileDetail.Uid,
						Uname: result.FileDetail.Uname,
						Ctim:  result.FileDetail.Ctim,
						Mtim:  result.FileDetail.Mtim,
						Atim:  result.FileDetail.Atim,
					},
					Content:     result.Content,
					Description: result.Description,
				},
			},
		}
		err := reportService.Client.Report(ReportEvent)
		if err != nil {
			log.Error(err)
			continue
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.AddCommand(report.MapReportCmd(cmd.MapImageCommand(scanImageCmd, scanImage), reportService))
	scanCmd.AddCommand(report.MapReportCmd(cmd.MapContainerCommand(scanContainerCmd, scanContainer), reportService))
	rootCmd.AddCommand(cmd.NewInfoCommand(pluginInfo))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
