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

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-log4j2/pkg/scanner"
)

var (
	pluginInfo = plugin.Manifest{
		Name:        "veinmind-log4j2",
		Author:      "veinmind-team",
		Description: "veinmind-log4j2 scanner image which has log4j jar vulnerable with CVE-2021-44228",
	}
	ReferencesURLList = []string{
		"https://www.kb.cert.org/vuls/id/930724",
		"https://nvd.nist.gov/vuln/detail/CVE-2021-44228",
		"https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2021-44228",
		"https://tools.cisco.com/security/center/content/CiscoSecurityAdvisory/cisco-sa-apache-log4j-qRuKNEbd",
	}
	reportService = &report.Service{}
	rootCmd       = &cmd.Command{}
	scanCmd       = &cmd.Command{
		Use: "scan",
	}
	scanImageCmd = &cmd.Command{
		Use:   "image",
		Short: "scan image command",
	}
	scanContainerCmd = &cmd.Command{
		Use:   "container",
		Short: "scan container command",
	}
)

func InitReferences() (res []event.References) {
	for _, value := range ReferencesURLList {
		tmpRef := event.References{
			Type: "URL",
			URL:  value,
		}
		res = append(res, tmpRef)
	}
	return res
}

func scanImage(c *cmd.Command, image api.Image) error {
	var result []*scanner.Result
	err := scanner.ScanImage(image, &result)
	if err != nil {
		log.Error("Scan Image Error")
		return err
	}
	if len(result) > 0 {
		for _, value := range result {
			reportEvent := &event.Event{
				BasicInfo: &event.BasicInfo{
					ID:         image.ID(),
					Object:     event.NewObject(image),
					Source:     pluginInfo.Name,
					Time:       time.Now(),
					Level:      event.Critical,
					DetectType: event.Image,
					EventType:  event.Risk,
					AlertType:  event.Vulnerability,
				},
				DetailInfo: event.NewDetailInfo(&event.VulnDetail{
					ID:         "CVE-2021-44228",
					Published:  time.Date(2021, 11, 26, 0, 0, 0, 0, time.Local),
					Summary:    "Apache Log4j2 2.0-beta9 through 2.15.0 (excluding security releases 2.12.2, 2.12.3, and 2.3.1) JNDI features used in configuration, log messages, and parameters do not protect against attacker controlled LDAP and other JNDI related endpoints.",
					Details:    "Apache Log4j2 2.0-beta9 through 2.15.0 (excluding security releases 2.12.2, 2.12.3, and 2.3.1) JNDI features used in configuration, log messages, and parameters do not protect against attacker controlled LDAP and other JNDI related endpoints. An attacker who can control log messages or log message parameters can execute arbitrary code loaded from LDAP servers when message lookup substitution is enabled. From log4j 2.15.0, this behavior has been disabled by default. From version 2.16.0 (along with 2.12.2, 2.12.3, and 2.3.1), this functionality has been completely removed. Note that this vulnerability is specific to log4j-core and does not affect log4net, log4cxx, or other Apache Logging Services projects.",
					References: InitReferences(),
					Source: event.Source{
						Type:     "jar",
						FilePath: value.DisplayPath,
						Packages: event.AssetPackageDetail{
							Name: value.File,
						},
					},
				}),
			}
			err := reportService.Client.Report(reportEvent)
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}

	return nil
}

func scanContainer(c *cmd.Command, container api.Container) error {
	var result []*scanner.Result
	err := scanner.ScanContainer(container, &result)
	if err != nil {
		log.Error("Scan Image Error")
		return err
	}
	if len(result) > 0 {
		for _, value := range result {
			reportEvent := &event.Event{
				BasicInfo: &event.BasicInfo{
					ID:         container.ID(),
					Object:     event.NewObject(container),
					Source:     pluginInfo.Name,
					Time:       time.Now(),
					Level:      event.Critical,
					DetectType: event.Container,
					EventType:  event.Risk,
					AlertType:  event.Vulnerability,
				},
				DetailInfo: event.NewDetailInfo(&event.VulnDetail{
					ID:         "CVE-2021-44228",
					Published:  time.Date(2021, 11, 26, 0, 0, 0, 0, time.Local),
					Summary:    "Apache Log4j2 2.0-beta9 through 2.15.0 (excluding security releases 2.12.2, 2.12.3, and 2.3.1) JNDI features used in configuration, log messages, and parameters do not protect against attacker controlled LDAP and other JNDI related endpoints.",
					Details:    "Apache Log4j2 2.0-beta9 through 2.15.0 (excluding security releases 2.12.2, 2.12.3, and 2.3.1) JNDI features used in configuration, log messages, and parameters do not protect against attacker controlled LDAP and other JNDI related endpoints. An attacker who can control log messages or log message parameters can execute arbitrary code loaded from LDAP servers when message lookup substitution is enabled. From log4j 2.15.0, this behavior has been disabled by default. From version 2.16.0 (along with 2.12.2, 2.12.3, and 2.3.1), this functionality has been completely removed. Note that this vulnerability is specific to log4j-core and does not affect log4net, log4cxx, or other Apache Logging Services projects.",
					References: InitReferences(),
					Source: event.Source{
						Type:     "jar",
						FilePath: value.DisplayPath,
						Packages: event.AssetPackageDetail{
							Name: value.File,
						},
					},
				}),
			}
			err := reportService.Client.Report(reportEvent)
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.AddCommand(report.MapReportCmd(cmd.MapContainerCommand(scanContainerCmd, scanContainer), reportService))
	scanCmd.AddCommand(report.MapReportCmd(cmd.MapImageCommand(scanImageCmd, scanImage), reportService))

	rootCmd.AddCommand(cmd.NewInfoCommand(pluginInfo))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
