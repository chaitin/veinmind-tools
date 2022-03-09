//+build community

package main

import (
	"github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	_ "github.com/chaitin/veinmind-tools/veinmind-malicious/config"
	_ "github.com/chaitin/veinmind-tools/veinmind-malicious/database"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/database/model"
	_ "github.com/chaitin/veinmind-tools/veinmind-malicious/database/model"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/scanner/malicious"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/sdk/common"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/sdk/common/report"
	"github.com/spf13/cobra"
	_ "net/http/pprof"
	"os"
	"strings"
	"sync"
	"time"
)

var reportData = model.ReportData{}
var reportLock sync.Mutex
var scanStart = time.Now()

var rootCmd = &cmd.Command{}

var scanCmd = &cmd.Command{
	Use: "scan",
	Short: "Scan image malicious files",
	PostRun: func(cmd *cobra.Command, args []string) {
		// 计算扫描数据
		spend := time.Since(scanStart)
		reportData.ScanSpendTime = spend.String()
		reportData.ScanStartTime = scanStart.Format("2006-01-02 15:04:05")
		report.CalculateScanReportCount(&reportData)
		report.SortScanReport(&reportData)

		format, _ := cmd.Flags().GetString("format")
		name, _ := cmd.Flags().GetString("name")
		name = strings.Join([]string{name, format}, ".")

		switch format {
		case report.HTML:
			report.OutputHTML(reportData, name)
		case report.JSON:
			report.OutputJSON(reportData, name)
		case report.CSV:
			report.OutputCSV(reportData, name)
		}
	},
}


func scan(_ *cmd.Command, image api.Image) error {
	report, err := malicious.Scan(image)
	if err != nil {
		common.Log.Error(err)
		return nil
	}

	reportLock.Lock()
	reportData.ScanImageResult = append(reportData.ScanImageResult, report)
	reportLock.Unlock()

	return nil
}

func init() {
	rootCmd.AddCommand(cmd.MapImageCommand(scanCmd, scan))
	rootCmd.AddCommand(cmd.NewInfoCommand(plugin.Manifest{}))
	scanCmd.Flags().StringP("format", "f", "html", "report format for scan report")
	scanCmd.Flags().StringP("name", "n", "report", "report name for scan report")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
