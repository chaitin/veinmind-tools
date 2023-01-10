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
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
)

var (
	Regix              = ""
	tempDir            = ""
	available_resource = []string{"pod", "pods", "configmap", "configmaps", "cm", "clusterrolebinding", "clusterrolebindings", "all"}
	scanIaCCmd         = &cmd.Command{
		Use:      "iac",
		Short:    "perform iac file scan",
		PreRunE:  scanReportPreRunE,
		RunE:     scanIaC,
		PostRunE: scanReportPostRunE,
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

func scanIaC(c *cmd.Command, args []string) error {
	wg := &sync.WaitGroup{}
	wg.Add(len(args))
	for _, value := range args {
		handler, err := ScanIaCParser(value)
		if err != nil {
			return err
		}
		go func(handler Handler, value string, wg *sync.WaitGroup) {
			defer wg.Done()
			err := handler(c, value)
			if err != nil {
				log.Errorf(err.Error())
			}
		}(handler, value, wg)
	}
	wg.Wait()
	return nil
}

// host
func scanHostIaCFile(c *cmd.Command, arg string) error {
	filetype, err := c.Flags().GetString("iac-type")
	var iacfile iacApi.IACType
	if err != nil {
		return err
	}
	regex := "(host:)?(.*)"
	compileRegex := regexp.MustCompile(regex)
	matchArr := compileRegex.FindStringSubmatch(arg)
	hostfilepath := matchArr[len(matchArr)-1]
	if filetype == "" {
		var opt []iacApi.DiscoverOption
		iacfile, err = iacApi.DiscoverType(hostfilepath, opt...)
		if err != nil {
			return err
		}
	} else {
		switch filetype {
		case "kubernetes":
			iacfile = "kubernetes"
		case "dockerfile":
			iacfile = "dockerfile"
		case "docker-compose":
			iacfile = "docker-compose"
		default:
			iacfile = "unknown"
		}
	}
	iac := iacApi.IAC{
		Path: hostfilepath,
		Type: iacfile,
	}
	_ = scanIaCFile(c, iac)
	return nil
}

// git
func scanGitRepoIaCFile(c *cmd.Command, arg string) error {
	key, err := c.Flags().GetString("sshkey")

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

	regex := "git:(.*)"
	compileRegex := regexp.MustCompile(regex)
	matchArr := compileRegex.FindStringSubmatch(arg)
	isGitUrl, err := regexp.MatchString("^(http(s)?://([^/]+?/){2}|git@[^:]+:[^/]+?/).*?.git$", matchArr[1])
	if err != nil {
		return err
	}
	if isGitUrl {
		func() {
			var opt []iacApi.DiscoverOption

			err = git.Clone(tempDir, matchArr[1], key, insecure)

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

	return nil
}

// kubernetes
func parserK8sConfig(c *cmd.Command, input string) ([]string, error) {
	var namespace, kubeconfig, resource, name string
	kubeconfig, _ = c.Flags().GetString("kubeconfig")
	if kubeconfig == "" {
		if home := os.Getenv("HOME"); home != "" {
			kubeconfig = home + "/.kube/config"
		} else {
			log.Errorf("please input kubeconfig file path using --kubeconfig / -k")
		}
	}
	if input != "kubernetes" {
		namespace, _ = c.Flags().GetString("namespace")
		regex := "(\\w*):(\\w*){1,}/?(.*)"
		compileRegex := regexp.MustCompile(regex)
		matchArr := compileRegex.FindStringSubmatch(input)
		resource = matchArr[2]
		name = matchArr[3]
	} else {
		namespace = "all"
		resource = "all"
		name = "all"
	}

	return []string{namespace, resource, name, kubeconfig}, nil
}
func inResource(input string) bool {
	for _, value := range available_resource {
		if strings.ToLower(input) == value {
			return true
		}
	}
	return false
}
func inNamespace(input string, availabe []string) bool {
	if input == "all" {
		return true
	} else {
		for _, value := range availabe {
			if input == value {
				return true
			}
		}
	}
	return false
}
func scanK8sConfig(c *cmd.Command, arg string) error {
	args, err := parserK8sConfig(c, arg)
	if err != nil {
		return err
	}
	namespace := args[0]
	resource := args[1]
	name := args[2]
	kubeconfig := args[3]
	if inResource(resource) == false {
		if resource == "" {
			log.Errorf("please input available resource!\n     available :pod,configmap,clusterrolebinding\n     get       :nil\nif you want to scan all resource please input\n    kubernetes or kubernetes:all/<yourname>")
		} else {
			log.Errorf("please input available resource!\n     available :pod,configmap,clusterrolebinding\n     get       :%s", resource)
		}

		return nil
	}
	log.Infof("start load remote k8s cluster config at %s", kubeconfig)
	scanCtx := c.Context()
	iacList := make([]iacApi.IAC, 0)
	dataList := map[string][]byte{}

	kubeRoot, err := kubernetes.New()
	if err != nil {
		return err
	}

	//配置namespace
	namespaces := make([]string, 0)
	k8sNamespaces, err := kubeRoot.Resource("", "namespaces")
	if err != nil {
		return err
	}
	available_namespace, err := k8sNamespaces.List(scanCtx)
	if err != nil {
		return err
	}
	if inNamespace(namespace, available_namespace) {
		if namespace == "all" {
			namespaces = available_namespace
		} else {
			namespaces = append(namespaces, namespace)
		}
	} else {
		log.Errorf("please input right namespace!\n    available namespaces:%s", available_namespace)
	}
	//配置name的正则表达式
	if name == "all" {
		Regix = ".*"
	} else {
		Regix = name
	}

	log.Infof("download remote k8s cluster config at %s", tempDir)
	resource = strings.ToLower(resource)
	for _, namespace := range namespaces {
		if kubeC, errName := kubernetes.New(); errName == nil {
			if resource == "all" || resource == "pod" || resource == "pods" {
				if resourcePod, errResource := kubeC.Resource(namespace, "pods"); errResource == nil {
					if pods, errPod := resourcePod.List(scanCtx); errPod == nil {
						for _, pod := range pods {
							if match, _ := regexp.MatchString(Regix, pod); match {
								if podsConfig, errConfig := resourcePod.Get(scanCtx, pod); errConfig == nil {
									dataList[strings.Join([]string{namespace, pod}, ":")] = podsConfig
								}
							}
						}
					}
				}
			}
			if resource == "all" || resource == "configmap" || resource == "configmaps" || resource == "cm" {
				if resourceCM, errResource := kubeC.Resource(namespace, "configmaps"); errResource == nil {
					if configmaps, errconfigmap := resourceCM.List(scanCtx); errconfigmap == nil {
						for _, configmap := range configmaps {
							if match, _ := regexp.MatchString(Regix, configmap); match {
								if configmapconfig, errConfigmaps := resourceCM.Get(scanCtx, configmap); errConfigmaps == nil {
									dataList[strings.Join([]string{namespace, configmap}, ":")] = configmapconfig
								}
							}

						}
					}
				}
			}
		}
	}
	if resource == "all" || resource == "clusterrolebinding" || resource == "clusterrolebindings" {
		if kubeC, errName := kubernetes.New(); errName == nil {
			if resourceClusterRole, errResource := kubeC.Resource("", "clusterrolebindings"); errResource == nil {
				if clusterRoleBindings, errclusterrolebinding := resourceClusterRole.List(scanCtx); errclusterrolebinding == nil {
					for _, clusterRoleBinding := range clusterRoleBindings {
						if match, _ := regexp.MatchString(name, clusterRoleBinding); match {
							if clusterRolebindingConfig, errConfig := resourceClusterRole.Get(scanCtx, clusterRoleBinding); errConfig == nil {
								dataList[strings.Join([]string{"none", clusterRoleBinding}, ":")] = clusterRolebindingConfig
							}
						}
					}
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

func ScanIaCParser(arg string) (Handler, error) {
	var flag string
	regex := "(kubernetes|host|git)?:?(.*)"
	compileRegex := regexp.MustCompile(regex)
	matchArr := compileRegex.FindStringSubmatch(arg)
	if matchArr[1] == "" { //没有协议头
		flag = "host"
	} else {
		flag = matchArr[1]
	}
	switch flag {
	case KUBERNETES:
		return scanK8sConfig, nil
	case GIT:
		return scanGitRepoIaCFile, nil
	case HOST:
		return scanHostIaCFile, nil
	}
	return nil, nil
}

func init() {
	scanIaCCmd.Flags().String("iac-type", "", "dedicate iac type for iac files")
	scanIaCCmd.Flags().StringP("kubeconfig", "k", "", "k8s config file")
	scanIaCCmd.Flags().StringP("namespace", "n", "all", "k8s resource namespace")
	scanIaCCmd.Flags().StringP("proxy", "", "", "proxy to git like: https://xxxxx or socks5://xxxx")
	scanIaCCmd.Flags().StringP("sshkey", "", "", "auth to git if use by ssh proto")
	scanIaCCmd.Flags().BoolP("insecure-skip", "", false, "skip tls config")
}
