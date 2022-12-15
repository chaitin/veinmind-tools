package main

import (
	"context"
	"github.com/chaitin/libveinmind/go/cmd"
	iacApi "github.com/chaitin/libveinmind/go/iac"
	"github.com/chaitin/libveinmind/go/kubernetes"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/libveinmind/go/plugin/service"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/git"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/plugind"
	"github.com/google/uuid"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	tempDir    = ""
	scanIaCCmd = &cmd.Command{
		Use:   "scan-iac",
		Short: "perform iac file scan",
	}
	scanLocalCmd = &cmd.Command{
		Use:      "local",
		Short:    "perform local iac file scan",
		PreRunE:  scanPreRunE,
		PostRunE: scanPostRunE,
	}
	scanGitRepoCmd = &cmd.Command{
		Use:      "git",
		Short:    "perform git repo iac file scan",
		PreRunE:  scanIacPreRunE,
		RunE:     scanGitRepoIaCFile,
		PostRunE: scanIacPostRunE,
	}
	scanK8sConfigCmd = &cmd.Command{
		Use:      "k8s",
		Short:    "perform scan iac by k8s config",
		PreRunE:  scanIacPreRunE,
		RunE:     scanK8sConfig,
		PostRunE: scanIacPostRunE,
	}

	// replace path
	scanIacPreRunE = func(c *cmd.Command, args []string) error {
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
	scanIacPostRunE = func(c *cmd.Command, args []string) error {
		if tempDir != "" {
			os.RemoveAll(tempDir)
		}
		return scanPostRunE(c, args)
	}
)

func scanIaCFile(c *cmd.Command, iac iacApi.IAC) error {
	log.Infof("scan Iac: %s, filetype: %s", iac.Path, iac.Type)
	if err := cmd.ScanIAC(ctx, ps, iac,
		plugin.WithExecInterceptor(func(
			ctx context.Context, plug *plugin.Plugin, c *plugin.Command, next func(context.Context, ...plugin.ExecOption) error,
		) error {
			// Register Service
			reg := service.NewRegistry()
			reg.AddServices(log.WithFields(log.Fields{
				"plugin":  plug.Name,
				"command": path.Join(c.Path...),
			}))
			reg.AddServices(reportService)

			// Next Plugin
			return next(ctx, reg.Bind())
		})); err != nil {
		return err
	}
	return nil
}

func scanGitRepoIaCFile(c *cmd.Command, args []string) error {
	key, err := c.Flags().GetString("ssh-pubkey")
	if err != nil {
		return err
	}

	insecure, err := c.Flags().GetBool("insecure-skip")
	if err != nil {
		return err
	}

	proxy, err := c.Flags().GetString("proxy")
	if err != nil {
		return err
	}

	if proxy != "" {
		os.Setenv("ALL_PROXY", proxy)
	}

	for _, arg := range args {
		isGitUrl, err := regexp.MatchString("^(http(s)?://([^/]+?/){2}|git@[^:]+:[^/]+?/).*?.git$", arg)
		if err != nil {
			continue
		}
		if isGitUrl {
			func() {
				var opt []iacApi.DiscoverOption
				err = git.Clone(tempDir, arg, key, insecure)
				if err != nil {
					log.Errorf("git download failed: %s", err)
					// nil point fix
					return
				}
				discovered, err := iacApi.DiscoverIACs(tempDir, opt...)
				if err != nil {
					log.Errorf("git discovered failed: %s", err)
				}
				for _, iac := range discovered {
					_ = scanIaCFile(c, iac)
				}
			}()
		}
	}

	return nil
}

func scanK8sConfig(c *cmd.Command, args []string) error {
	kubeconfig, _ := c.Flags().GetString("kubeconfig")
	log.Infof("start load remote k8s cluster config at %s", kubeconfig)
	scanCtx := c.Context()
	iacList := make([]iacApi.IAC, 0)
	dataList := map[string][]byte{}

	option := kubernetes.WithKubeConfig(kubeconfig)
	kubeRoot, err := kubernetes.New(option)

	if err != nil {
		return err
	}
	kubeNamespaces, err := kubeRoot.Resource("namespaces")
	if err != nil {
		return err
	}

	namespaces, err := kubeNamespaces.List(scanCtx)
	if err != nil {
		return err
	}

	log.Infof("download remote k8s cluster config at %s", tempDir)
	for _, namespace := range namespaces {
		optionForNamespace := kubernetes.WithNamespace(namespace)
		if kubeC, errName := kubernetes.New(option, optionForNamespace); errName == nil {
			// pods
			if resourcePod, errResource := kubeC.Resource("pods"); errResource == nil {
				if pods, errPod := resourcePod.List(scanCtx); errPod == nil {
					for _, pod := range pods {
						if podsConfig, errConfig := resourcePod.Get(scanCtx, pod); errConfig == nil {
							dataList[strings.Join([]string{namespace, pod}, ":")] = podsConfig
						}
					}
				}
			}
			// configMaps
			if resourceConfigMaps, errResource := kubeC.Resource("configmaps"); errResource == nil {
				if configmaps, errCM := resourceConfigMaps.List(scanCtx); errCM == nil {
					for _, configmap := range configmaps {
						if kubeletconfig, errConfig := resourceConfigMaps.Get(scanCtx, configmap); errConfig == nil {
							dataList[strings.Join([]string{namespace, configmap}, ":")] = kubeletconfig
						}
					}
				}
			}
		}
	}

	// some Others
	if resourceClusterRole, errResource := kubeRoot.Resource("clusterrolebindings"); errResource == nil {
		if clusterRoleBindings, errClusterRolebinding := resourceClusterRole.List(scanCtx); errClusterRolebinding == nil {
			for _, clusterRoleBinding := range clusterRoleBindings {
				if clusterRolebindingConfig, errConfig := resourceClusterRole.Get(scanCtx, clusterRoleBinding); errConfig == nil {
					dataList[strings.Join([]string{"none", clusterRoleBinding}, ":")] = clusterRolebindingConfig
				}
			}
		}
	}

	// write TempFile
	for key, data := range dataList {
		tmpConfigFile := path.Join(tempDir, key)
		if errWrite := ioutil.WriteFile(tmpConfigFile, data, fs.ModePerm); errWrite == nil {
			iacList = append(iacList, iacApi.IAC{
				Path: tmpConfigFile,
				Type: iacApi.Kubernetes,
			})
		} else {
			log.Warnf("write temp file failed: %s", tmpConfigFile)
		}
	}

	for _, iac := range iacList {
		_ = scanIaCFile(c, iac)
	}

	return nil
}

func init() {
	scanIaCCmd.AddCommand(cmd.MapIACCommand(scanLocalCmd, scanIaCFile))
	scanIaCCmd.AddCommand(scanGitRepoCmd)
	scanIaCCmd.AddCommand(scanK8sConfigCmd)

	scanIaCCmd.PersistentFlags().String("iac-type", "", "dedicate iac type for iac files")
	scanIaCCmd.PersistentFlags().Int("threads", 5, "threads for scan action")
	scanIaCCmd.PersistentFlags().StringP("output", "o", "report.json", "output filepath of report")
	scanIaCCmd.PersistentFlags().StringP("glob", "g", "", "specifies the pattern of plugin file to find")

	scanGitRepoCmd.Flags().String("proxy", "", "proxy to git like: https://xxxxx or socks5://xxxx")
	scanGitRepoCmd.Flags().String("sshkey", "", "auth to git if use by ssh proto")
	scanGitRepoCmd.Flags().Bool("insecure-skip", false, "skip tls config")

	scanK8sConfigCmd.Flags().String("kubeconfig", "", "k8s config file")
}
