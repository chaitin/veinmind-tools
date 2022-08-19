package main

import (
	"context"
	"io/fs"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	reportService "github.com/chaitin/veinmind-common-go/service/report"
	"golang.org/x/sync/errgroup"

	"veinmind-webshell/pkg/detect"
	"veinmind-webshell/pkg/filter"
)

var (
	rootCommand = &cmd.Command{}
	scanCommand = &cmd.Command{
		Use:   "scan image webshell",
		Short: "scan image webshell",
	}

	// flags
	token string
)

func scan(c *cmd.Command, image api.Image) (err error) {
	detectKit, err := detect.NewKit(context.Background(), detect.WithToken(token), detect.WithDefaultClient())
	if err != nil {
		return err
	}

	// Error group for event handler
	reportEvents := make(chan reportService.ReportEvent, 1<<8)
	errGHandler := errgroup.Group{}
	errGHandler.Go(func() error {
		for evt := range reportEvents {
			err := reportService.DefaultReportClient().Report(evt)
			if err != nil {
				log.Error(err)
				continue
			}
		}
		return nil
	})

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
				reportEvents <- *evt

				return nil
			})

			return nil
		}
	})
	errG.Wait()
	close(reportEvents)
	errGHandler.Wait()

	return err
}

func init() {
	rootCommand.AddCommand(cmd.MapImageCommand(scanCommand, scan))
	rootCommand.AddCommand(cmd.NewInfoCommand(plugin.Manifest{
		Name:        "veinmind-webshell",
		Author:      "veinmind-team",
		Description: "veinmind-webshell scan image webshell data",
	}))
	scanCommand.Flags().StringVarP(&token, "token", "t", "", "百川 api token")
}

func main() {
	if err := rootCommand.ExecuteContext(context.Background()); err != nil {
		panic(err)
	}
}
