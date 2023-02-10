package main

import (
	"context"
	"io/fs"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-webshell/pkg/detect"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-webshell/pkg/filter"
	"golang.org/x/sync/errgroup"
)

var (
	reportService = &report.Service{}
	rootCommand   = &cmd.Command{}

	scanCommand = &cmd.Command{
		Use:   "scan",
		Short: "scan mode",
	}
	scanImageCommand = &cmd.Command{
		Use:   "image",
		Short: "scan image basic info",
	}
	scanContainerCommand = &cmd.Command{
		Use:   "container",
		Short: "scan container basic info",
	}

	// flags
	token string
)

func scanImage(c *cmd.Command, image api.Image) (err error) {
	detectKit, err := detect.NewKit(context.Background(), detect.WithToken(token), detect.WithDefaultClient())
	if err != nil {
		return err
	}

	// Error group for detect kit
	errG := errgroup.Group{}
	errG.SetLimit(100)

	err = image.Walk("/", func(path string, info fs.FileInfo, err error) error {
		if isScript, scriptType, err := filter.Kit.Filter(path, info); !isScript {
			return nil
		} else {
			if err != nil {
				log.Error(err)
				return nil
			}

			errG.Go(func() error {
				f, err := image.Open(path)
				if err != nil {
					log.Error(err)
					return nil
				}

				detectFileInfo := detect.FileInfo{
					Path:        path,
					Reader:      f,
					ScriptType:  scriptType,
					RawFileInfo: info,
				}
				res, err := detectKit.Detect(detectFileInfo)
				if err != nil {
					log.Error(err)
					return nil
				}

				// Send result to channel
				evt, err := detect.Convert2ReportEvent(image, detectFileInfo, *res)
				if err != nil {
					log.Error(err)
					return nil
				}
				if evt != nil {
					err := reportService.Client.Report(evt)
					if err != nil {
						return err
					}
				}
				return nil
			})

			return nil
		}
	})
	errG.Wait()

	return err
}

func scanContainer(c *cmd.Command, container api.Container) (err error) {
	detectKit, err := detect.NewKit(context.Background(), detect.WithToken(token), detect.WithDefaultClient())
	if err != nil {
		return err
	}

	// Error group for detect kit
	errG := errgroup.Group{}
	errG.SetLimit(100)

	err = container.Walk("/", func(path string, info fs.FileInfo, err error) error {
		if isScript, scriptType, err := filter.Kit.Filter(path, info); !isScript {
			return nil
		} else {
			if err != nil {
				log.Error(err)
				return nil
			}

			errG.Go(func() error {
				f, err := container.Open(path)
				if err != nil {
					log.Error(err)
					return nil
				}

				detectFileInfo := detect.FileInfo{
					Path:        path,
					Reader:      f,
					ScriptType:  scriptType,
					RawFileInfo: info,
				}
				res, err := detectKit.Detect(detectFileInfo)
				if err != nil {
					log.Error(err)
					return nil
				}

				// Send result to channel
				evt, err := detect.Convert2ReportEvent(container, detectFileInfo, *res)
				if err != nil {
					log.Error(err)
					return nil
				}
				if evt != nil {
					err := reportService.Client.Report(evt)
					if err != nil {
						return err
					}
				}
				return nil
			})

			return nil
		}
	})
	errG.Wait()

	return err
}

func init() {
	rootCommand.AddCommand(scanCommand)
	rootCommand.AddCommand(cmd.NewInfoCommand(plugin.Manifest{
		Name:        "veinmind-webshell",
		Author:      "veinmind-team",
		Description: "veinmind-webshell scan image webshell data",
	}))

	scanCommand.AddCommand(report.MapReportCmd(cmd.MapImageCommand(scanImageCommand, scanImage), reportService))
	scanCommand.AddCommand(report.MapReportCmd(cmd.MapContainerCommand(scanContainerCommand, scanContainer), reportService))
	scanCommand.PersistentFlags().StringVarP(&token, "token", "t", "", "百川 api token")
}

func main() {
	if err := rootCommand.ExecuteContext(context.Background()); err != nil {
		panic(err)
	}
}
