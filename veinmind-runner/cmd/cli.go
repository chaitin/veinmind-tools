package main

import (
	"context"
	_ "embed"
	"errors"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"

	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/container"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/plugind"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/reporter"
)

const (
	resourceDirectoryPath = "./resource"
)

var (
	ps                    []*plugin.Plugin
	ctx                   context.Context
	allPlugins            []*plugin.Plugin //All found plugins
	serviceManager        *plugind.Manager
	cancel                context.CancelFunc
	runnerReporter        *reporter.Reporter
	reportService         *report.ReportService
	parallelContainerMode = container.InContainer()
	scanPreRunE           = func(c *cobra.Command, args []string) error {
		// create resource directory if not exist
		if _, err := os.Open(resourceDirectoryPath); os.IsNotExist(err) {
			err = os.Mkdir(resourceDirectoryPath, 0600)
			if err != nil {
				return err
			}
		}

		// discover plugins
		ctx = c.Context()
		ctx, cancel = context.WithCancel(ctx)
		ps = []*plugin.Plugin{}

		glob, err := c.Flags().GetString("glob")
		if err == nil && glob != "" {
			allPlugins, err = plugin.DiscoverPlugins(ctx, ".", plugin.WithGlob(glob))
		} else {
			allPlugins, err = plugin.DiscoverPlugins(ctx, ".")
		}
		if err != nil {
			return err
		}

		serviceManager, err = plugind.NewManager()
		if err != nil {
			return err
		}

		for _, p := range allPlugins {
			log.Infof("Discovered plugin: %#v\n", p.Name)
			err = serviceManager.StartWithContext(ctx, p.Name)
			if err != nil {
				log.Errorf("%#v can not work: %#v\n", p.Name, err)
				continue
			}
			ps = append(ps, p)
		}

		// reporter channel listen
		go runnerReporter.Listen()

		// event channel listen
		go func() {
			for {
				select {
				case evt := <-reportService.EventChannel:
					runnerReporter.EventChannel <- evt
				}
			}
		}()

		return nil
	}
	scanPostRunE = func(cmd *cobra.Command, args []string) error {
		// Stop reporter listen
		runnerReporter.StopListen()

		// Output
		err := runnerReporter.Write(os.Stdout)
		if err != nil {
			log.Error(err)
		}
		output, _ := cmd.Flags().GetString("output")
		if parallelContainerMode {
			output = filepath.Join(resourceDirectoryPath, output)
		}
		if _, err := os.Stat(output); errors.Is(err, os.ErrNotExist) {
			f, err := os.Create(output)
			if err != nil {
				log.Error(err)
			} else {
				err = runnerReporter.Write(f)
				if err != nil {
					return err
				}
			}
		} else {
			f, err := os.OpenFile(output, os.O_WRONLY, 0666)
			if err != nil {
				log.Error(err)
			} else {
				err = runnerReporter.Write(f)
				if err != nil {
					return err
				}
			}
		}

		cancel()
		serviceManager.Wait()

		// Exit
		exitcode, err := cmd.Flags().GetInt("exit-code")
		if err != nil {
			return err
		}

		if exitcode == 0 {
			return nil
		} else {
			events, err := runnerReporter.GetEvents()
			if err != nil {
				return err
			}

			if len(events) > 0 {
				os.Exit(exitcode)
			} else {
				return nil
			}
		}

		return nil
	}
)

// root command
var rootCmd = &cmd.Command{}

func init() {
	// Cobra init
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(scanHostCmd)
	rootCmd.AddCommand(scanIaCCmd)
	rootCmd.AddCommand(scanRegistryCmd)

	// control exit
	rootCmd.PersistentFlags().IntP("exit-code", "e", 0, "exit-code when veinmind-runner find security issues")

	// Service client init
	reportService = report.NewReportService()

	// Reporter init
	r, err := reporter.NewReporter()
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
