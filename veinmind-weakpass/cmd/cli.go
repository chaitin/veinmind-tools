package main

import (
	"fmt"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/veinmind-common/go/service/report"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/embed"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/model"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/scanner"
	"github.com/spf13/cobra"
	_ "net/http/pprof"
	"os"
	"strconv"
	"sync"
	"text/tabwriter"
	"time"
)

var results = []model.ScanImageResult{}
var resultsLock sync.Mutex
var scanStart = time.Now()

var rootCmd = &cmd.Command{}
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "extract dict file to disk",
	Run: func(cmd *cobra.Command, args []string) {
		embed.ExtractAll()
	},
}
var scanCmd = &cmd.Command{
	Use:   "scan",
	Short: "Scan image weakpass",
	PostRun: func(cmd *cobra.Command, args []string) {
		tabw := tabwriter.NewWriter(os.Stdout, 95, 95, 0, ' ', tabwriter.TabIndent|tabwriter.Debug)
		fmt.Fprintln(tabw, "# ============================================================================================ #")
		spend := time.Since(scanStart)
		fmt.Fprintln(tabw, "| Scan Total: ", strconv.Itoa(len(results)), "\t")
		fmt.Fprintln(tabw, "| Spend Time: ", spend.String(), "\t")
		var weakpassImageTotal = 0
		var weakpassTotal = 0
		for _, r := range results {
			if len(r.WeakpassResults) > 0 {
				weakpassImageTotal++
				weakpassTotal += len(r.WeakpassResults)
			}
		}
		fmt.Fprintln(tabw, "| Weakpass Image Total: ", strconv.Itoa(weakpassTotal), "\t")
		fmt.Fprintln(tabw, "| Weakpass Total: ", strconv.Itoa(weakpassTotal), "\t")
		fmt.Fprintln(tabw, "+----------------------------------------------------------------------------------------------+")

		for _, r := range results {
			if len(r.WeakpassResults) > 0 {
				fmt.Fprintln(tabw, "| ImageName: ", r.ImageName, "\t")
				fmt.Fprintln(tabw, "| Status: Unsafe", "\t")
				for _, w := range r.WeakpassResults {
					fmt.Fprintln(tabw, "| Username: ", w.Username, "\t")
					fmt.Fprintln(tabw, "| Password: ", w.Password, "\t")
				}
				fmt.Fprintln(tabw, "+----------------------------------------------------------------------------------------------+")
			}
		}
		fmt.Fprintln(tabw, "# ============================================================================================ #\n")
		tabw.Flush()
	},
}

func scan(c *cmd.Command, image api.Image) error {
	result, err := scanner.Scan(image, scanner.ScanOption{
		ScanThreads: func() int {
			threads, err := c.Flags().GetInt("threads")
			if err != nil {
				return 10
			} else {
				return threads
			}
		}(),
		Username: func() string {
			username, err := c.Flags().GetString("username")
			if err != nil {
				return ""
			} else {
				return username
			}
		}(),
		Dictpath: func() string {
			dictpath, err := c.Flags().GetString("dictpath")
			if err != nil {
				return ""
			} else {
				return dictpath
			}
		}(),
	})

	if err != nil {
		log.Error(err)
		return nil
	}

	resultsLock.Lock()
	results = append(results, result)
	resultsLock.Unlock()

	// Report Event
	if len(result.WeakpassResults) > 0 {
		details := []report.AlertDetail{}
		for _, wr := range result.WeakpassResults {
			details = append(details, report.AlertDetail{
				WeakpassDetail: &report.WeakpassDetail{
					Username: wr.Username,
					Password: wr.Password,
					Service:  report.WeakpassService(wr.PassType)},
			})
		}
		reportEvent := report.ReportEvent{
			ID:           image.ID(),
			Time:         time.Now(),
			Level:        report.High,
			DetectType:   report.Image,
			EventType:    report.Risk,
			AlertType:    report.Weakpass,
			AlertDetails: details,
		}
		err = report.DefaultReportClient().Report(reportEvent)
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
		Name:        "veinmind-weakpass",
		Author:      "veinmind-team",
		Description: "veinmind-weakpass scanner image weakpass",
	}))
	scanCmd.Flags().IntP("threads", "t", 10, "password brute threads")
	scanCmd.Flags().StringP("username", "u", "", "username e.g. root")
	scanCmd.Flags().StringP("dictpath", "d", "", "dict path e.g ./mypass.dict")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
