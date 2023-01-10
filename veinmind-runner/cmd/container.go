package main

import (
	"context"
	"fmt"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/containerd"
	"github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/libveinmind/go/plugin/service"
	"path"
	"regexp"
	"sync"
)

var scanContainerCmd = &cmd.Command{
	Use:      "container",
	Short:    "perform container scan ",
	PreRunE:  scanReportPreRunE,
	RunE:     ScanContainer,
	PostRunE: scanReportPostRunE,
}

func ScanContainer(c *cmd.Command, args []string) error {
	if len(args) == 0 {
		args = append(args, "")
	}
	wg := &sync.WaitGroup{}
	wg.Add(len(args))
	for _, value := range args {
		handler, err := ScanContainerParser(value)
		if err != nil {
			return err
		}
		go func(handler Handler, value string, wg *sync.WaitGroup) {
			defer wg.Done()
			err := handler(c, value)
			if err != nil {
				log.Error(err)
			}
		}(handler, value, wg)
	}
	wg.Wait()
	return nil
}

func ScanContainerParser(arg string) (Handler, error) {
	compileRegex := regexp.MustCompile(CONTAINERREGEX)
	matchArr := compileRegex.FindStringSubmatch(arg)
	switch matchArr[1] {
	case ALL: //如果没指定运行时则先以输入的id/name作为docker尝试打开容器，若无法打开则为containerd容器 否则为docker容器
		r, err := docker.New()
		if err != nil {
			log.Error(err)
			return nil, err
		}
		ids := make([]string, 0)
		ids, err = r.FindContainerIDs(arg)
		if arg == "" {
			ids, err = r.ListContainerIDs()
		}
		if err != nil {
			log.Error(err)
			return nil, err
		}
		for _, id := range ids {
			_, err := r.OpenContainerByID(id)
			if err != nil {
				return ScanContainerd, nil
			}
		}
		return ScanDocker, nil
	case DOCKER:
		return ScanDocker, nil
	case CONTAINERD:
		return ScanContainerd, nil
	default:
		log.Errorf("please input right args! available: docker , containerd")
	}
	return nil, nil
}

func ScanDocker(c *cmd.Command, arg string) error {
	compileRegex := regexp.MustCompile(DOCKERREGEX)
	matchArr := compileRegex.FindStringSubmatch(arg)
	ids := make([]string, 0)
	r, err := docker.New()
	if err != nil {
		log.Error(err)
		return err
	}
	if matchArr[len(matchArr)-1] == "" {
		refs, err := r.ListContainerIDs()
		if err != nil {
			log.Error(err)
			return err
		}
		for _, ref := range refs {
			tmp, err := r.FindContainerIDs(ref)
			if err != nil {
				log.Error(err)
				continue
			}
			ids = append(ids, tmp...)
		}
	} else {
		ids, err = r.FindContainerIDs(matchArr[len(matchArr)-1])
		if err != nil {
			log.Error(err)
			return err
		}
	}
	for _, id := range ids {
		runtime, err := r.OpenContainerByID(id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		containerScan(c, runtime)
	}
	return nil

}

func ScanContainerd(c *cmd.Command, arg string) error {
	compileRegex := regexp.MustCompile(CONTAINERDREGEX)
	matchArr := compileRegex.FindStringSubmatch(arg)
	ids := make([]string, 0)
	r, err := containerd.New()
	if err != nil {
		return err
	}
	if matchArr[len(matchArr)-1] == "" {
		refs, err := r.ListContainerIDs()
		if err != nil {
			log.Error(err)
			return err
		}
		for _, ref := range refs {
			tmp, err := r.FindContainerIDs(ref)
			if err != nil {
				log.Error(err)
				continue
			}
			ids = append(ids, tmp...)
		}
	} else {
		ids, err = r.FindContainerIDs(matchArr[len(matchArr)-1])
		if err != nil {
			log.Error(err)
			return err
		}
	}
	for _, id := range ids {
		runtime, err := r.OpenContainerByID(id)
		if err != nil {
			return err
		}
		containerScan(c, runtime)
	}

	return nil
}

func containerScan(c *cmd.Command, container api.Container) error {

	ref := container.Name()

	// Get threads value
	t, err := c.Flags().GetInt("threads")
	if err != nil {
		t = 5
	}

	log.Infof("Scan container: %#v\n", ref)
	if err := cmd.ScanContainer(ctx, ps, container,
		plugin.WithExecInterceptor(func(
			ctx context.Context, plug *plugin.Plugin, c *plugin.Command, next func(context.Context, ...plugin.ExecOption) error,
		) error {
			// Register Service
			reg := service.NewRegistry()
			opts := make([]plugin.ExecOption, 0)
			opts = append(opts, reg.Bind())
			if value, ok := pluginArgsMap[plug.Name]; ok == true {
				opts = append(opts, plugin.WithPrependArgs(value...))
			}
			reg.AddServices(log.WithFields(log.Fields{
				"plugin":  plug.Name,
				"command": path.Join(c.Path...),
			}))
			reg.AddServices(reportService)

			// Next Plugin
			return next(ctx, opts...)
		}), plugin.WithExecParallelism(t)); err != nil {
		return err
	}
	return nil
}
