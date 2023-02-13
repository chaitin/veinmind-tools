package main

import (
	"context"
	_ "embed"
	"os"

	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report"

	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/reporter"
)

const (
	resourceDirectoryPath = "./resource"
)

// root command
var rootCmd = &cmd.Command{
	Use:  "veinmind-runner",
	Long: `veinmind-runner is a veinmind container security tool platform developed by Chaitin Technology`,
}

func init() {
	// with context
	rootCmd.SetContext(context.Background())

	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(generateCmd)

	// global params
	rootCmd.PersistentFlags().Int("threads", 5, "threads for scan action")
	rootCmd.PersistentFlags().StringP("output", "o", "report.json", "output filepath of report")
	rootCmd.PersistentFlags().StringP("glob", "g", "", "specifies the pattern of plugin file to find")
	rootCmd.PersistentFlags().IntP("exit-code", "e", 0, "exit-code when veinmind-runner find security issues")
	// Service client init
	reportService = report.NewReportService()

	// Reporter init
	r, err := reporter.NewReporter(rootCmd.Context())
	if err != nil {
		log.Fatal(err)
	}
	runnerReporter = r
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
