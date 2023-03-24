package main

import (
	"context"
	_ "embed"
	"os"

	"github.com/chaitin/libveinmind/go/cmd"
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

	rootCmd.AddCommand(aiCmd)
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(generateCmd)

	// global params
	rootCmd.PersistentFlags().Int("threads", 5, "threads for scan action")
	rootCmd.PersistentFlags().StringP("output", "o", "", "output filepath of report")
	rootCmd.PersistentFlags().StringP("glob", "g", "", "specifies the pattern of plugin file to find")
	rootCmd.PersistentFlags().IntP("exit-code", "e", 0, "exit-code when veinmind-runner find security issues")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
