package main

import (
	"os"
	"time"

	"github.com/chaitin/libveinmind/go/cmd"
	iacApi "github.com/chaitin/libveinmind/go/iac"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-common-go/service/report/event"
	"github.com/open-policy-agent/opa/ast"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-iac/pkg/scanner"
)

var (
	pluginInfo = plugin.Manifest{
		Name:        "veinmind-iac",
		Author:      "veinmind-team",
		Description: "veinmind-iac scan IAC file and discovery risks of them",
	}

	reportService = &report.Service{}

	results   []scanner.Result
	scanTotal = 0

	reportLevelMap = map[string]event.Level{
		"Low":      event.Low,
		"Medium":   event.Medium,
		"High":     event.High,
		"Critical": event.Critical,
	}
	scanCmd = &cmd.Command{
		Use: "scan",
	}
	rootCmd    = &cmd.Command{}
	scanIaCCmd = &cmd.Command{
		Use:   "iac",
		Short: "scan iac command",
	}
)

func scanIaC(c *cmd.Command, iac iacApi.IAC) error {
	scanTotal += 1
	// do something here
	scanner := &scanner.Scanner{
		QueryPre: "data.brightMirror.",
		Policies: make(map[string]*ast.Module),
	}
	scanner.LoadLibs()
	res, err := scanner.Scan(c.Context(), iac)
	if err != nil {
		log.Error(err)
		return err
	}

	uniqueAppend(res)

	for _, data := range res {
		for _, risk := range data.Risks {
			reportEvent := event.Event{
				BasicInfo: &event.BasicInfo{
					ID:         iac.Path,
					Object:     event.NewObject(iac),
					Time:       time.Now(),
					Source:     pluginInfo.Name,
					Level:      reportLevelMap[data.Rule.Severity],
					DetectType: event.Container,
					EventType:  event.Risk,
					AlertType:  event.IaCRisk,
				},
				DetailInfo: event.NewDetailInfo(&event.IaCDetail{
					RuleInfo: event.IaCRule{
						Id:          data.Rule.Id,
						Name:        data.Rule.Name,
						Description: data.Rule.Description,
						Reference:   data.Rule.Reference,
						Severity:    data.Rule.Severity,
						Solution:    data.Rule.Solution,
						Type:        data.Rule.Type,
					},
					FileInfo: event.IaCData{
						StartLine: risk.StartLine,
						EndLine:   risk.EndLine,
						FilePath:  risk.FilePath,
						Original:  risk.Original,
					},
				}),
			}
			err = reportService.Client.Report(&reportEvent)
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
	scanCmd.AddCommand(report.MapReportCmd(cmd.MapIACCommand(scanIaCCmd, scanIaC), reportService))
	rootCmd.AddCommand(cmd.NewInfoCommand(pluginInfo))
}

func uniqueAppend(res []scanner.Result) {
	for _, tr := range res {
		uniqueFlag := false
		for i, r := range results {
			if r.Id == tr.Id {
				uniqueFlag = true
				results[i].Risks = append(results[i].Risks, tr.Risks...)
				break
			}
		}
		if !uniqueFlag {
			results = append(results, tr)
		}
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
