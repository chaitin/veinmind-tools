package main

import (
	"os"
	"time"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-common-go/service/report/event"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-minio/pkg/scanner"
)

var reportService = &report.Service{}
var rootCmd = &cmd.Command{}

var scanCmd = &cmd.Command{
	Use: "scan",
}

var scanImageCmd = &cmd.Command{
	Use:   "image",
	Short: "scan image command",
}

var scanContainerCmd = &cmd.Command{
	Use:   "container",
	Short: "scan container command",
}

var ReferencesURLList = []string{
	"https://mp.weixin.qq.com/s/JgskenAZ6Cpecoe2k2AEjQ",
	"https://github.com/minio/minio/releases/tag/RELEASE.2023-03-20T20-16-18Z",
	"https://github.com/minio/minio/security/advisories/GHSA-6xvq-wj2x-3h3q",
}

// scanImage is func that used to do some action with Images
// you can write your plugin scan code here
func scanImage(c *cmd.Command, image api.Image) error {
	// do something here
	res := scanner.ScanImage(image)

	if res.Version != "" {
		// if you want display at runner report, you should send your result to report event
		reportEvent := &event.Event{
			BasicInfo: &event.BasicInfo{
				ID:         image.ID(),
				Object:     event.NewObject(image),
				Time:       time.Now(),
				Level:      event.Critical,
				DetectType: event.Image,
				EventType:  event.Risk,
				AlertType:  event.Vulnerability,
			},
			DetailInfo: &event.DetailInfo{
				AlertDetail: &event.VulnDetail{
					ID:         "CVE-2023-28432",
					Published:  time.Date(2023, 03, 23, 0, 0, 0, 0, time.Local),
					Summary:    "Information Disclosure in Cluster Deployment",
					Details:    "In a cluster deployment, MinIO returns all environment variables, including MINIO_SECRET_KEY and MINIO_ROOT_PASSWORD, resulting in information disclosure. All users of distributed deployment are impacted. All users are advised to upgrade ASAP.",
					References: initReferences(),
					Source: event.Source{
						Type:     "go-binary",
						FilePath: res.File,
						Packages: event.AssetPackageDetail{
							Name:    res.File,
							Version: res.Version,
						},
					},
				},
			},
		}
		err := reportService.Client.Report(reportEvent)
		if err != nil {
			return err
		}
	}

	return nil
}

// scanContainer is func that used to do some action with container
// you can write your plugin scan code here
func scanContainer(c *cmd.Command, container api.Container) error {
	// do something here
	res := scanner.ScanContainer(container)

	if res.Version != "" {
		reportEvent := &event.Event{
			BasicInfo: &event.BasicInfo{
				ID:         container.ID(), // container id info
				Object:     event.NewObject(container),
				Time:       time.Now(),      // report time, usually use time.Now
				Level:      event.Critical,  // report event level
				DetectType: event.Container, // report scan object type
				EventType:  event.Risk,
				AlertType:  event.Vulnerability,
			},
			DetailInfo: &event.DetailInfo{
				AlertDetail: &event.VulnDetail{
					ID:         "CVE-2023-28432",
					Published:  time.Date(2023, 03, 23, 0, 0, 0, 0, time.Local),
					Summary:    "Information Disclosure in Cluster Deployment",
					Details:    "In a cluster deployment, MinIO returns all environment variables, including MINIO_SECRET_KEY and MINIO_ROOT_PASSWORD, resulting in information disclosure. All users of distributed deployment are impacted. All users are advised to upgrade ASAP.",
					References: initReferences(),
					Source: event.Source{
						Type:     "go-binary",
						FilePath: res.File,
						Packages: event.AssetPackageDetail{
							Name:    res.File,
							Version: res.Version,
						},
					},
				},
			},
		}
		err := reportService.Client.Report(reportEvent)
		if err != nil {
			return err
		}
	}

	return nil
}

func initReferences() (res []event.References) {
	for _, value := range ReferencesURLList {
		tmpRef := event.References{
			Type: "URL",
			URL:  value,
		}
		res = append(res, tmpRef)
	}
	return res
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.AddCommand(report.MapReportCmd(cmd.MapImageCommand(scanImageCmd, scanImage), reportService))
	scanCmd.AddCommand(report.MapReportCmd(cmd.MapContainerCommand(scanContainerCmd, scanContainer), reportService))
	rootCmd.AddCommand(cmd.NewInfoCommand(plugin.Manifest{
		Name:        "veinmind-minio",
		Author:      "veinmind-team",
		Description: "veinmind-minio scan CVE-2023-28432 risk in images/containers",
	}))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
