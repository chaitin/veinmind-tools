package main

import (
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-iac/pkg/output"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-iac/pkg/scanner"
	"github.com/open-policy-agent/opa/ast"
	"os"
	"time"

	"github.com/chaitin/libveinmind/go/cmd"
	iacApi "github.com/chaitin/libveinmind/go/iac"
	"github.com/chaitin/libveinmind/go/plugin"
)

var (
	results   []scanner.Result
	scanStart = time.Now()
	scanTotal = 0

	reportLevelMap = map[string]report.Level{
		"Low":      report.Low,
		"Medium":   report.Medium,
		"High":     report.High,
		"Critical": report.Critical,
	}
	rootCmd    = &cmd.Command{}
	scanIaCCmd = &cmd.Command{
		Use:   "scan-iac",
		Short: "scan image command",
		PostRun: func(cmd *cmd.Command, args []string) {
			format, _ := cmd.Flags().GetString("format")
			spend := time.Since(scanStart)

			output.Stream(spend, scanTotal, func() error {
				switch format {
				case output.STDOUT:
					if err := output.Stdout(results); err != nil {
						log.Error("Stdout error", err)
						return err
					}
				case output.JSON:
					if err := output.Json(results); err != nil {
						log.Error("Export Results JSON False", err)
						return err
					}
				}
				return nil
			})
		},
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

	reportDetails := make([]report.AlertDetail, 0)
	reportLevel := report.Low
	for _, data := range res {
		for _, risk := range data.Risks {
			if reportLevel < reportLevelMap[data.Rule.Severity] {
				reportLevel = reportLevelMap[data.Rule.Severity]
			}
			reportDetails = append(reportDetails, report.AlertDetail{
				IaCDetail: &report.IaCDetail{
					RuleInfo: report.IaCRule{
						Id:          data.Rule.Id,
						Name:        data.Rule.Name,
						Description: data.Rule.Description,
						Reference:   data.Rule.Reference,
						Severity:    data.Rule.Severity,
						Solution:    data.Rule.Solution,
						Type:        data.Rule.Type,
					},
					FileInfo: report.IaCData{
						StartLine: risk.StartLine,
						EndLine:   risk.EndLine,
						FilePath:  risk.FilePath,
						Original:  risk.Original,
					},
				},
			})
		}
	}

	if len(reportDetails) > 0 {
		//if you want display at runner report, you should send your result to report event
		reportEvent := report.ReportEvent{
			ID:             iac.Path,                 // image id info
			Time:           time.Now(),               // report time, usually use time.Now
			Level:          reportLevel,              // report event level
			DetectType:     report.IaC,               // report scan object type
			EventType:      report.Risk,              // report event type: Risk/Invasion/Info
			AlertType:      report.IaCRisk,           // report alert type, we provide some clearly types of security events,
			AlertDetails:   reportDetails,            // add report detail data in there
			GeneralDetails: []report.GeneralDetail{}, // if your report event does not in alert type, you can use GeneralDetails type which consists of json bytes
		}
		err = report.DefaultReportClient(report.WithDisableLog()).Report(reportEvent)
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(cmd.MapIACCommand(scanIaCCmd, scanIaC))
	rootCmd.AddCommand(cmd.NewInfoCommand(plugin.Manifest{
		Name:        "veinmind-iac",
		Author:      "veinmind-team",
		Description: "veinmind-iac scan IAC file and discovery risks of them",
	}))
	scanIaCCmd.Flags().StringP("format", "f", "stdout", "export file format")
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
