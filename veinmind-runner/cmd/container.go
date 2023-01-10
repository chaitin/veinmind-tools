package main

import (
	"context"
	"fmt"
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
	regex := "(docker|containerd)?:?(.*)"
	compileRegex := regexp.MustCompile(regex)
	matchArr := compileRegex.FindStringSubmatch(arg)
	if matchArr[1] == "" || matchArr[1] == "docker" {
		return ScanDocker, nil
	} else if matchArr[1] == "containerd" {
		return ScanContainerd, nil
	} else {
		log.Errorf("please input right args! available: docker , containerd")
	}
	return nil, nil
}

func ScanDocker(c *cmd.Command, arg string) error {
	regex := "docker?:?(.*)"
	compileRegex := regexp.MustCompile(regex)
	matchArr := compileRegex.FindStringSubmatch(arg)
	r, err := docker.New()
	if err != nil {
		return err
	}
	ids, err := r.FindContainerIDs(matchArr[1])
	for _, id := range ids {
		runtime, err := r.OpenContainerByID(id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		ref := runtime.Name()

		// Get threads value
		t, err := c.Flags().GetInt("threads")
		if err != nil {
			t = 5
		}

		log.Infof("Scan container: %#v\n", ref)
		if err := cmd.ScanContainer(c.Context(), ps, runtime,
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
			}), plugin.WithExecParallelism(t)); err != nil {
			return err
		}
	}
	return nil

}

func ScanContainerd(c *cmd.Command, arg string) error {
	regex := "containerd:?(.*)"
	compileRegex := regexp.MustCompile(regex)
	matchArr := compileRegex.FindStringSubmatch(arg)
	r, err := containerd.New()
	if err != nil {
		return err
	}
	ids, err := r.FindContainerIDs(matchArr[1])
	for _, id := range ids {
		container, err := r.OpenContainerByID(id)
		if err != nil {
			return err
		}
		ref := container.Name()
		return nil
		// Get threads value
		t, err := c.Flags().GetInt("threads")
		if err != nil {
			t = 5
		}

		log.Infof("Scan container: %#v\n", ref)
		if err := cmd.ScanContainer(c.Context(), ps, container,
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
			}), plugin.WithExecParallelism(t)); err != nil {
			return err
		}
	}

	return nil
}
