package scan

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/chaitin/libveinmind/go/cmd"
	iacApi "github.com/chaitin/libveinmind/go/iac"
	"github.com/chaitin/libveinmind/go/kubernetes"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/gogf/gf/errors/gerror"
	"golang.org/x/sync/errgroup"

	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/log"

	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/git"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/target"
)

var AvailableResource = []string{"pod", "pods", "configmap", "configmaps", "cm", "clusterrolebinding", "clusterrolebindings", "all"}

func DispatchIacs(ctx context.Context, targets []*target.Target) error {
	errG := errgroup.Group{}
	for _, obj := range targets {
		errG.Go(func() error {
			switch obj.Proto {
			case target.LOCAL:
				return HostIac(ctx, obj)
			case target.GIT:
				return GitIac(ctx, obj)
			case target.KUBERNETES:
				return KubeIac(ctx, obj)
			default:
				return errors.New(fmt.Sprintf("individual iac proto: %s", obj.Proto))
			}
		})
	}
	return errG.Wait()
}

func HostIac(ctx context.Context, t *target.Target) error {
	// 如果用户没有输入扫描路径，代表从当前路径开始自动扫描。
	if t.Value == "" {
		t.Value = "./"
	}
	// check dir / file
	info, err := os.Stat(t.Value)
	if err != nil {
		return gerror.Wrap(err, "un exits file info")
	}

	// dir scan
	if info.IsDir() {
		var iacOpt []iacApi.DiscoverOption
		if t.Opts.IacLimitSize > 0 {
			iacOpt = append(iacOpt, iacApi.WithIACLimitSize(t.Opts.IacLimitSize))
		}
		if t.Opts.IacFileType != "" {
			iacOpt = append(iacOpt, iacApi.WithIACType(iacApi.IACType(t.Opts.IacFileType)))
		}
		iacFiles, err := iacApi.DiscoverIACs(t.Value, iacOpt...)
		if err != nil {
			return gerror.Wrap(err, "auto discover iac error")
		}
		for _, file := range iacFiles {
			if err := doIAC(ctx, t.Plugins, file, t.WithDefaultOptions()...); err != nil {
				log.GetModule(log.ScanModuleKey).Errorf("scan iac %s error: %+v", file.Path, err)
			}
		}
		return nil
	}
	var fileType iacApi.IACType
	if !iacApi.IsIACType(t.Opts.IacFileType) {
		fileType, err = iacApi.DiscoverType(t.Value)
		if err != nil {
			return err
		}
	} else {
		fileType = iacApi.IACType(t.Opts.IacFileType)
	}
	return doIAC(ctx, t.Plugins, iacApi.IAC{
		Path: t.Value,
		Type: fileType,
	}, t.WithDefaultOptions()...)
}

func GitIac(ctx context.Context, t *target.Target) error {
	if t.Opts.IacProxy != "" {
		os.Setenv("ALL_PROXY", t.Opts.IacProxy)
	}
	isGitUrl, err := regexp.MatchString("^(http(s)?://([^/]+?/){2}|git@[^:]+:[^/]+?/).*?.git$", t.Value)
	if err != nil {
		return err
	}
	if isGitUrl {
		err = git.Clone(t.Opts.TempPath, t.Value, t.Opts.IacSshPath, t.Opts.Insecure)

		if err != nil {
			log.GetModule(log.ScanModuleKey).Errorf("git download failed: %+v", err)
			// nil point fix
			return err
		}
		return HostIac(ctx, &target.Target{
			Proto:          t.Proto,
			Value:          t.Opts.TempPath,
			Opts:           t.Opts,
			Plugins:        t.Plugins,
			ReportService:  t.ReportService,
			ServiceManager: t.ServiceManager,
		})
	}
	return errors.New("not git url")
}

func KubeIac(ctx context.Context, t *target.Target) error {
	args, err := parserK8sConfig(t)
	if err != nil {
		return err
	}
	namespace := args[0]
	resource := args[1]
	name := args[2]
	kubeconfig := args[3]
	if inResource(resource) == false {
		if resource == "" {
			log.GetModule(log.ScanModuleKey).Errorf("please input available resource!\n     available :pod,configmap,clusterrolebinding\n     get       :nil\nif you want to scan all resource please input\n    kubernetes or kubernetes:all/<yourname>")
		} else {
			log.GetModule(log.ScanModuleKey).Errorf("please input available resource!\n     available :pod,configmap,clusterrolebinding\n     get       :%s", resource)
		}

		return nil
	}
	log.GetModule(log.ScanModuleKey).Infof("start load remote k8s cluster config at %s", kubeconfig)
	dataList := map[string][]byte{}
	kubeRoot, err := kubernetes.New(kubernetes.WithKubeConfigPath(kubeconfig))
	if err != nil {
		return err
	}

	//配置namespace
	namespaces := make([]string, 0)
	k8sNamespaces, err := kubeRoot.Resource("", "namespaces")
	if err != nil {
		return err
	}

	available_namespace, err := k8sNamespaces.List(ctx)
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
		log.GetModule(log.ScanModuleKey).Errorf("please input right namespace!\n    available namespaces:%s", available_namespace)
	}
	reg := name
	//配置name的正则表达式
	if reg == "all" {
		reg = ".*"
	}

	log.GetModule(log.ScanModuleKey).Infof("download remote k8s cluster config at %s", t.Opts.TempPath)
	resource = strings.ToLower(resource)
	for _, namespace := range namespaces {
		if resource == "all" || resource == "pod" || resource == "pods" {
			if resourcePod, errResource := kubeRoot.Resource(namespace, "pods"); errResource == nil {
				if pods, errPod := resourcePod.List(ctx); errPod == nil {
					for _, pod := range pods {
						if match, _ := regexp.MatchString(reg, pod); match {
							if podsConfig, errConfig := resourcePod.Get(ctx, pod); errConfig == nil {
								dataList[strings.Join([]string{namespace, pod}, ":")] = podsConfig
							}
						}
					}
				}
			}
		}
		if resource == "all" || resource == "configmap" || resource == "configmaps" || resource == "cm" {
			if resourceCM, errResource := kubeRoot.Resource(namespace, "configmaps"); errResource == nil {
				if configmaps, errconfigmap := resourceCM.List(ctx); errconfigmap == nil {
					for _, configmap := range configmaps {
						if match, _ := regexp.MatchString(reg, configmap); match {
							if configmapconfig, errConfigmaps := resourceCM.Get(ctx, configmap); errConfigmaps == nil {
								dataList[strings.Join([]string{namespace, configmap}, ":")] = configmapconfig
							}
						}

					}
				}
			}
		}
	}
	if resource == "all" || resource == "clusterrolebinding" || resource == "clusterrolebindings" {
		if resourceClusterRole, errResource := kubeRoot.Resource("", "clusterrolebindings"); errResource == nil {
			if clusterRoleBindings, errclusterrolebinding := resourceClusterRole.List(ctx); errclusterrolebinding == nil {
				for _, clusterRoleBinding := range clusterRoleBindings {
					if match, _ := regexp.MatchString(name, clusterRoleBinding); match {
						if clusterRolebindingConfig, errConfig := resourceClusterRole.Get(ctx, clusterRoleBinding); errConfig == nil {
							dataList[strings.Join([]string{"none", clusterRoleBinding}, ":")] = clusterRolebindingConfig
						}
					}
				}
			}
		}
	}
	// write TempFile
	for key, data := range dataList {
		tmpConfigFile := path.Join(t.Opts.TempPath, key)
		if errWrite := ioutil.WriteFile(tmpConfigFile+".yaml", data, fs.ModePerm); errWrite == nil {
			log.GetModule(log.ScanModuleKey).Infof("write temp config file at %s", tmpConfigFile)
		} else {
			log.GetModule(log.ScanModuleKey).Warnf("write temp file failed: %s", tmpConfigFile)
		}
	}

	return HostIac(ctx, &target.Target{
		Proto:          t.Proto,
		Value:          t.Opts.TempPath,
		Opts:           t.Opts,
		Plugins:        t.Plugins,
		ReportService:  t.ReportService,
		ServiceManager: t.ServiceManager,
	})
}

func parserK8sConfig(t *target.Target) ([]string, error) {
	var namespace, resource, name, kubeconfig string
	kubeconfig = t.Opts.IacKubeConfig
	if kubeconfig == "" {
		if home := os.Getenv("HOME"); home != "" {
			kubeconfig = home + "/.kube/config"
		} else {
			log.GetModule(log.ScanModuleKey).Errorf("please input kubeconfig file path using --kubeconfig / -k")
		}
	}
	if t.Value == "" {
		namespace = "all"
		resource = "all"
		name = "all"
	} else {
		namespace = t.Opts.IacKubeNameSpace
		if namespace == "" {
			namespace = "all"
		}
		target := strings.Split(t.Value, "/")
		if len(target) == 1 {
			resource = target[0]
			name = "all"
		} else if len(target) == 2 {
			resource = target[0]
			name = target[1]
		}
	}
	return []string{namespace, resource, name, kubeconfig}, nil
}

func inResource(input string) bool {
	for _, value := range AvailableResource {
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

func doIAC(ctx context.Context, rang plugin.ExecRange, iac iacApi.IAC, pluginOpts ...plugin.ExecOption) error {
	log.GetModule(log.ScanModuleKey).Infof("start scan iac: %s, filetype: %s", iac.Path, iac.Type)
	return cmd.ScanIAC(ctx, rang, iac, pluginOpts...)
}
