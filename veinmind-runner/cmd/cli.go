package main

import (
	"context"
	_ "embed"
	"errors"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"

	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/container"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/plugind"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/reporter"
)

const (
	CONTAINERREGEX  string = "(docker|containerd)?:?(.*)"
	IMAGEREGEX      string = "(docker|containerd|registry)?:?(.*)"
	IACREGEX        string = "(kubernetes|host|git)?:?(.*)"
	DOCKERREGEX     string = "(docker)?:?(.*)"
	CONTAINERDREGEX string = "(containerd)?:?(.*)"
	REGISTRYREGEX   string = "(registry)?:?(.*)"
	KUBERNETESREGEX string = "(\\w*):(\\w*){1,}/?(.*)"
	GITREGEX        string = "(git)?:?(.*)"
	HOSTREGEX       string = "(host)?:?(.*)"

	resourceDirectoryPath = "./resource"
	ALL                   = ""
	DOCKER                = "docker"
	CONTAINERD            = "containerd"
	REGISTRY              = "registry"
	KUBERNETES            = "kubernetes"
	GIT                   = "git"
	HOST                  = "host"
)

var (
	tempDir               = ""
	ps                    []*plugin.Plugin
	ctx                   context.Context
	allPlugins            []*plugin.Plugin //All found plugins
	serviceManager        *plugind.Manager
	cancel                context.CancelFunc
	runnerReporter        *reporter.Reporter
	reportService         *report.ReportService
	parallelContainerMode = container.InContainer()
	pluginArgs            = make(map[string][]string, 0)
	pluginArgsMap         = make(map[string][]string, 0)
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
			ps, err = plugin.DiscoverPlugins(ctx, ".", plugin.WithGlob(glob))
		} else {
			ps, err = plugin.DiscoverPlugins(ctx, ".")
		}
		if err != nil {
			return err
		}
		pluginArgsMap = parseArgs(c, args, ps)
		serviceManager, err = plugind.NewManager()
		if err != nil {
			return err
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

	scanReportPreRunE = func(c *cmd.Command, args []string) error {
		// init tempDir
		dir, err := os.MkdirTemp("", uuid.NewString())
		if err != nil {
			return err
		}
		tempDir = dir
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
					// replace temp path
					if realID, err := filepath.Rel(tempDir, evt.ID); err == nil {
						evt.ID = realID
					}
					for _, detail := range evt.AlertDetails {
						if detail.IaCDetail != nil {
							if realPath, err := filepath.Rel(tempDir, detail.IaCDetail.FileInfo.FilePath); err == nil {
								detail.IaCDetail.FileInfo.FilePath = realPath
							}
						}
					}
					runnerReporter.EventChannel <- evt
				}
			}
		}()

		return nil
	}
	scanReportPostRunE = func(c *cmd.Command, args []string) error {
		if tempDir != "" {
			os.Remove(tempDir)
		}
		return scanPostRunE(c, args)
	}
)

type Handler func(c *cmd.Command, arg string) error

func splitArgs(cmd *cobra.Command, args []string) ([]string, []string) {
	if cmd.ArgsLenAtDash() >= 0 {
		return args[:cmd.ArgsLenAtDash()], args[cmd.ArgsLenAtDash():]
	}
	return args, []string{}
}

func parseArgs(cmd *cmd.Command, args []string, ps []*plugin.Plugin) map[string][]string {
	_, targetArgs := splitArgs(cmd, args)
	targetArgs = append(targetArgs, "--")
	res := make(map[string][]string, 0)
	for _, plugin := range ps {
		res[plugin.Name] = []string{}
	}
	for i := 0; i < len(targetArgs); i++ {
		if targetArgs[i] == "--" {
			tmp := targetArgs[0:i]
			if _, ok := res[tmp[0]]; ok == true {
				res[tmp[0]] = append(res[tmp[0]], tmp...)
			}
			targetArgs = targetArgs[i+1:]
			i = 0
		}
	}
	return res
}

// root command
var rootCmd = &cmd.Command{}

func init() {

	// Cobra init
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(scanIaCCmd)
	rootCmd.AddCommand(scanContainerCmd)
	rootCmd.AddCommand(scanImageCmd)

	rootCmd.PersistentFlags().Int("threads", 5, "threads for scan action")
	rootCmd.PersistentFlags().StringP("output", "o", "report.json", "output filepath of report")
	rootCmd.PersistentFlags().StringP("glob", "g", "", "specifies the pattern of plugin file to find")
	// control exit
	rootCmd.PersistentFlags().IntP("exit-code", "e", 0, "exit-code when veinmind-runner find security issues")

	// Scan Flags
	scanImageCmd.Flags().StringP("config", "c", "", "auth config path")
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
