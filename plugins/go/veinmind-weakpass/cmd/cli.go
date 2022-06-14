package main

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/dict/embed"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/model"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/utils"

	"github.com/spf13/cobra"
)

var results = []model.ScanImageResult{}
var serviceName = []string{}
var threads int
var username string
var dictpath string
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
		fmt.Fprintln(tabw, "| Weakpass Image Total: ", strconv.Itoa(weakpassImageTotal), "\t")
		fmt.Fprintln(tabw, "| Weakpass Total: ", strconv.Itoa(weakpassTotal), "\t")
		fmt.Fprintln(tabw, "+----------------------------------------------------------------------------------------------+")

		for _, r := range results {
			if len(r.WeakpassResults) > 0 {
				fmt.Fprintln(tabw, "| ImageName: ", r.ImageName, "\t")
				fmt.Fprintln(tabw, "| ServiceName: ", r.ServiceName, "\t")
				fmt.Fprintln(tabw, "| Status: Unsafe", "\t")
				for _, w := range r.WeakpassResults {
					fmt.Fprintln(tabw, "| Username: ", w.Username, "\t")
					fmt.Fprintln(tabw, "| Password: ", w.Password, "\t")
					fmt.Fprintln(tabw, "| Filepath: ", w.Filepath, "\t")
				}
				fmt.Fprintln(tabw, "+----------------------------------------------------------------------------------------------+")
			}
		}
		fmt.Fprintln(tabw, "# ============================================================================================ #")
		tabw.Flush()
	},
}

func scan(c *cmd.Command, image api.Image) (err error) {
	config := model.Config{Thread: threads, Username: username, Dictpath: dictpath}
	for _, service := range serviceName {
		ModuleResult, err := utils.StartModule(config, image, service)
		if err != nil {
			log.Error(err)
			continue
		}
		results = append(results, ModuleResult)
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
	scanCmd.Flags().IntVarP(&threads, "threads", "t", 10, "password brute threads")
	scanCmd.Flags().StringVarP(&username, "username", "u", "", "username e.g. root")
	scanCmd.Flags().StringVarP(&dictpath, "dictpath", "d", "", "dict path e.g. ./mypass.dict")
	scanCmd.Flags().StringSliceVarP(&serviceName, "serviceName", "s", []string{"mysql", "tomcat", "redis", "ssh"}, "find weakpass in these service e.g. ssh")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
