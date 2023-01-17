package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/container"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/plugind"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/reporter"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/scan"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/target"
	"github.com/google/uuid"
)

// scan cmd
// scan support: image, container, iac
var (
	tempDir               = ""
	ps                    []*plugin.Plugin
	ctx                   context.Context
	serviceManager        *plugind.Manager
	cancel                context.CancelFunc
	runnerReporter        *reporter.Reporter
	reportService         *report.ReportService
	parallelContainerMode = container.InContainer()

	scanCmd = &cmd.Command{
		Use:   "scan",
		Short: "perform a scan",
	}
	scanImageCmd = &cmd.Command{
		Use:      "image",
		Short:    "perform image scan",
		PreRunE:  scanPreRun,
		PostRunE: scanPostRun,
	}
	scanContainerCmd = &cmd.Command{
		Use:      "container",
		Short:    "perform container scan",
		PreRunE:  scanPreRun,
		PostRunE: scanPostRun,
	}
	scanIaCCmd = &cmd.Command{
		Use:      "iac",
		Short:    "perform iac file scan",
		PreRunE:  scanPreRun,
		PostRunE: scanPostRun,
	}
)

func MapTaskCommand(c *cmd.Command, t scan.DispatchTask) *cmd.Command {
	c.RunE = func(c *cmd.Command, args []string) error {
		opts := generateOptions(c, args)
		// trans users args to scan target
		objs := target.NewTargets(c, args, ps, serviceManager, reportService, opts...)
		return t(ctx, objs)
	}
	return c
}

func generateOptions(c *cmd.Command, args []string) []target.Option {
	// trans users flags to target options
	opts := make([]target.Option, 0)
	// threads
	thread, _ := c.Flags().GetInt("threads")
	opts = append(opts, target.WithThread(thread))
	// insecure
	insecure, _ := c.Flags().GetBool("insecure-skip")
	opts = append(opts, target.WithInsecure(insecure))
	// config
	config, _ := c.Flags().GetString("config")
	opts = append(opts, target.WithConfigPath(config))
	// parallelMode
	opts = append(opts, target.WithParallelContainerMode(parallelContainerMode))
	// tempDir
	opts = append(opts, target.WithTempPath(tempDir))
	// resourceDir
	opts = append(opts, target.WithResourcePath(resourceDirectoryPath))
	// Iac param
	if c.Use == "iac" {
		// Iac Type
		iacType, _ := c.Flags().GetString("iac-type")
		opts = append(opts, target.WithIacFileType(iacType))
		// proxy
		iacProxy, _ := c.Flags().GetString("proxy")
		opts = append(opts, target.WithIacProxy(iacProxy))
		// sshkey
		iacSsh, _ := c.Flags().GetString("sshkey")
		opts = append(opts, target.WithIacSshPath(iacSsh))
		// kubeconfig
		iacConfig, _ := c.Flags().GetString("kubeconfig")
		opts = append(opts, target.WithIacKubeConfig(iacConfig))
		// kube namespace
		iacNamespace, _ := c.Flags().GetString("namespace")
		opts = append(opts, target.WithIacKubeNameSpace(iacNamespace))
	}
	return opts
}

func scanPreRun(c *cmd.Command, args []string) error {
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
		ps, err = plugin.DiscoverPlugins(ctx, ".", plugin.WithGlob(glob))
	} else {
		ps, err = plugin.DiscoverPlugins(ctx, ".")
	}
	if err != nil {
		return err
	}

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
				// iac need replace real path
				if c.Use == "iac" {
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
				}
				runnerReporter.EventChannel <- evt
			}
		}
	}()

	return nil
}

func scanPostRun(c *cmd.Command, args []string) error {
	if tempDir != "" {
		err := os.RemoveAll(tempDir)
		if err != nil {
			log.Errorf("remove temp dir err: %s", err)
		}
	}
	// Stop reporter listen
	runnerReporter.StopListen()

	// Output
	err := runnerReporter.Write(os.Stdout)
	if err != nil {
		log.Error(err)
	}
	output, _ := c.Flags().GetString("output")
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
	exitcode, err := c.Flags().GetInt("exit-code")
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

func init() {
	scanCmd.AddCommand(MapTaskCommand(scanImageCmd, scan.DispatchImages))
	scanCmd.AddCommand(MapTaskCommand(scanContainerCmd, scan.DispatchContainers))
	scanCmd.AddCommand(MapTaskCommand(scanIaCCmd, scan.DispatchIacs))
	scanCmd.PersistentFlags().BoolP("insecure-skip", "", false, "skip tls config")
	// Scan Flags
	scanImageCmd.Flags().StringP("config", "c", "", "auth config path")

	scanIaCCmd.Flags().String("iac-type", "", "dedicate iac type for iac files")
	scanIaCCmd.Flags().StringP("proxy", "", "", "proxy to git like: https://xxxxx or socks5://xxxx")
	scanIaCCmd.Flags().StringP("sshkey", "", "", "auth to git if use by ssh proto")
	scanIaCCmd.Flags().StringP("kubeconfig", "k", "", "k8s config file")
	scanIaCCmd.Flags().StringP("namespace", "n", "all", "k8s resource namespace")
}
