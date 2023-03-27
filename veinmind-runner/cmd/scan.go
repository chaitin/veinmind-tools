package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-common-go/service/report/service"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/container"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/log"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/plugind"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/scan"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/target"
)

// scan cmd
// scan support: image, container, iac
var (
	tempDir               = ""
	ps                    []*plugin.Plugin
	ctx                   context.Context
	serviceManager        *plugind.Manager
	cancel                context.CancelFunc
	reportService         *report.Service
	parallelContainerMode = container.InContainer()

	scanCmd = &cmd.Command{
		Use:   "scan",
		Short: "Scan cloud native objects security, include image/container/iac",
	}
	scanImageCmd = &cmd.Command{
		Use:   "image [flags] target",
		Short: `Scan image`,
		Long:  `Scan image from multi source, include dockerd/containerd/registry`,
		Example: `
1. scan dockerd image nginx:latest
veinmind-runner scan image dockerd:nginx:latest

2. scan containerd image bitnami/nginx:latest
veinmind-runner scan image containerd:bitnami/nginx:latest

3. scan public registry image library/ubuntu (all tag)
veinmind-runner scan image registry-image:library/ubuntu

4. scan private registry image example.com/app/market:v1.11.2
veinmind-runner scan image -c auth.toml registry-image:example.com/app/market:v1.11.2

5. use openai analyze result
veinmind-runner scan image --enable-analyze --openai-token  <your_openai_token> nginx:latest

6. use openai analyze result with yourself questions
veinmind-runner scan image --enable-analyze --openai-token  <your_openai_token> -p "explain what happened at this json" nginx:latest

7. scan private registry example.com （need admin privilege)
veinmind-runner scan image -c auth.toml registry:example.com

auth.toml format (yaml):
[[auths]]
	registry = "example.com"
	username = "<your-username>"
	password = "<your-password>"

8. scan tarball format image
veinmind-runner scan image tarball:/tmp/alpine.tar
`,
		PreRunE:  scanPreRun,
		PostRunE: scanPostRun,
	}
	scanContainerCmd = &cmd.Command{
		Use:   "container [flags] target",
		Short: "Scan container",
		Long:  `Scan container from multi source, include dockerd/containerd`,
		Example: `
1. scan dockerd container (all)
veinmind-runner scan container dockerd:*

2. scan dockerd container d29e2ca5b3a8 (container id)
veinmind-runner scan container dockerd:d29e2ca5b3a8

3. scan containerd container webapp (container name)
veinmind-runner scan container containerd:webapp

4. use openai analyze result
veinmind-runner scan container --enable-analyze --openai-token <your_openai_token> containerd:webapp

5. use openai analyze result with yourself questions
veinmind-runner scan container --enable-analyze --openai-token <your_openai_token> -p "explain what happened at this json" containerd:webapp
`,
		PreRunE:  scanPreRun,
		PostRunE: scanPostRun,
	}
	scanIaCCmd = &cmd.Command{
		Use:   "iac [flags] target",
		Short: "Scan iac",
		Long:  `Scan iac from multi source, include host/git/kubernetes`,
		Example: `
1. scan host iac (current directory)
veinmind-runner scan iac host:./

2. scan github kubernetes-sigs/kustomize repo iac
veinmind-runner scan iac git:https://github.com/kubernetes-sigs/kustomize.git

3. scan kubernetes pod iac
veinmind-runner scan iac --kubeconfig admin.yaml kubernetes:pod

4. use openai analyze result
veinmind-runner scan iac --kubeconfig admin.yaml --enable-analyze --openai-token  <your_openai_token> kubernetes:pod

5. use openai analyze result with yourself questions
veinmind-runner scan iac --kubeconfig admin.yaml --enable-analyze --openai-token  <your_openai_token> -p "explain what happened at this json" kubernetes:pod
`,
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
	// filter regex
	regex, _ := c.Flags().GetString("filter")
	opts = append(opts, target.WithFilterRegex(regex))
	// parallelMode
	opts = append(opts, target.WithParallelContainerMode(parallelContainerMode))
	// tempDir
	opts = append(opts, target.WithTempPath(tempDir))
	// resourceDir
	opts = append(opts, target.WithResourcePath(resourceDirectoryPath))
	// Iac param
	if c.Name() == "iac" {
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

	output, _ := c.Flags().GetString("output")
	if parallelContainerMode {
		output = filepath.Join(resourceDirectoryPath, output)
	}

	opts := make([]service.Option, 0)

	opts = append(opts, service.WithOutputDir(output))

	format, _ := c.Flags().GetString("format")
	formatList := strings.Split(format, ",")
	for _, f := range formatList {
		switch f {
		case "cli":
			opts = append(opts, service.WithTableRender())
		case "json":
			opts = append(opts, service.WithJsonRender())
		case "html":
			opts = append(opts, service.WithHtmlRender())
		}
	}

	verbose, _ := c.Flags().GetBool("verbose")
	if verbose {
		opts = append(opts, service.WithVerbose())
	}
	// discover plugins
	ctx = c.Context()
	// Service client init
	reportService = report.NewService(ctx, opts...)
	ctx, cancel = context.WithCancel(ctx)
	ps = []*plugin.Plugin{}

	glob, err := c.Flags().GetString("glob")
	if err == nil && glob != "" {
		ps, err = plugin.DiscoverPlugins(ctx, "./plugin", plugin.WithGlob(glob))
	} else {
		ps, err = plugin.DiscoverPlugins(ctx, "./plugin")
	}
	if err != nil {
		return err
	}

	serviceManager, err = plugind.NewManager()
	if err != nil {
		return err
	}

	// reporter channel listen
	go reportService.Listen()

	return nil
}

func scanPostRun(c *cmd.Command, _ []string) error {
	if tempDir != "" {
		err := os.RemoveAll(tempDir)
		if err != nil {
			log.GetModule(log.ScanModuleKey).Errorf(errors.Wrap(err, "can't remove temp directory").Error())
		}
	}
	// Stop reporter listen
	reportService.Close()
	// AI analyze
	analyze, err := c.Flags().GetBool("enable-analyze")
	if analyze {
		log.GetModule(log.ScanModuleKey).Infof("enbale ai analyzer, prepare use openai to analyze results......")
		token, _ := c.Flags().GetString("openai-token")
		if token != "" {
			prefix, _ := c.Flags().GetString("prefix")
			err := AnalyzeReport(ctx, token, prefix, reportService.EventPool.Events)
			if err != nil {
				log.GetModule(log.ScanModuleKey).Errorf(errors.Wrap(err, "openai analyze error").Error())
			}
		} else {
			log.GetModule(log.ScanModuleKey).Errorf(errors.New("empty openai_key, if you want use openai analyze results, use `-t/--token` with your openai_key").Error())
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
		if len(reportService.EventPool.Events) > 0 {
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

	scanCmd.PersistentFlags().StringP("format", "", "cli", "output format: html/json/cli")
	scanCmd.PersistentFlags().BoolP("verbose", "v", false, "output verbose detail")
	scanCmd.PersistentFlags().Bool("enable-analyze", false, "auto use openai analyze result")
	scanCmd.PersistentFlags().StringP("openai-token", "", "", "auto openai analyze openai_key")
	scanCmd.PersistentFlags().StringP("prefix", "p", "", "training openai limit sentence")
	scanCmd.PersistentFlags().BoolP("insecure-skip", "", false, "skip tls config")
	// Scan Flags
	scanImageCmd.Flags().StringP("config", "c", "", "auth config path")
	scanImageCmd.Flags().StringP("filter", "f", "", "catalog repo filter regex")

	scanIaCCmd.Flags().String("iac-type", "", "dedicate iac type for iac files")
	scanIaCCmd.Flags().StringP("proxy", "", "", "proxy to git like: https://xxxxx or socks5://xxxx")
	scanIaCCmd.Flags().StringP("sshkey", "", "", "auth to git if use by ssh proto")
	scanIaCCmd.Flags().StringP("kubeconfig", "k", "", "k8s config file")
	scanIaCCmd.Flags().StringP("namespace", "n", "all", "k8s resource namespace")
}
