package scan

import (
	"context"
	"errors"
	"fmt"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/containerd"
	"github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/libveinmind/go/plugin"
	"golang.org/x/sync/errgroup"

	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/log"

	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/target"
)

func DispatchContainers(ctx context.Context, targets []*target.Target) error {
	errG := errgroup.Group{}
	for _, obj := range targets {
		errG.Go(func() error {
			switch obj.Proto {
			case target.DOCKERD:
				r, err := docker.New()
				if err != nil {
					return err
				}
				return LocalContainer(ctx, obj, r)
			case target.CONTAINERD:
				r, err := containerd.New()
				if err != nil {
					return err
				}
				return LocalContainer(ctx, obj, r)
			default:
				return errors.New(fmt.Sprintf("[scan] individual container proto: %s", obj.Proto))
			}
		})
	}
	return errG.Wait()
}

func LocalContainer(ctx context.Context, t *target.Target, runtime api.Runtime) error {
	var ids []string
	var err error
	// scanAll
	if t.Value == "" || t.Value == "*" {
		ids, err = runtime.ListContainerIDs()
	} else {
		ids, err = runtime.FindContainerIDs(t.Value)
	}
	if err != nil {
		return err
	}
	for _, id := range ids {
		container, err := runtime.OpenContainerByID(id)
		if err != nil {
			return err
		}
		if err := doContainer(ctx, t.Plugins, container, t.WithDefaultOptions()...); err != nil {
			log.GetModule(log.ScanModuleKey).Errorf("[scan] scan container %s error : %s", container.Name(), err)
		}
	}
	return nil
}

func doContainer(ctx context.Context, rang plugin.ExecRange, container api.Container, pluginOpts ...plugin.ExecOption) error {
	log.GetModule(log.ScanModuleKey).Infof("[scan] start scan container: %#v\n", container.Name())
	return cmd.ScanContainer(ctx, rang, container, pluginOpts...)
}
