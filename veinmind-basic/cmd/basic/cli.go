package main

import (
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/veinmind-common/go/service/report"
	"os"
	"time"
)

var (
	rootCommand = &cmd.Command{}
	scanCommand = &cmd.Command{
		Use:   "scan",
		Short: "scan image basic info",
	}
)

func scan(c *cmd.Command, image api.Image) error {
	refs, err := image.RepoRefs()
	if err != nil {
		// no reference image will report ans use sha256 fill repo field
		log.Error(err)
	}

	oci, err := image.OCISpecV1()
	if err != nil {
		return err
	}

	evt := report.ReportEvent{
		ID:         image.ID(),
		Time:       time.Now(),
		Level:      report.None,
		DetectType: report.Image,
		AlertType:  report.Basic,
		EventType:  report.Info,
		AlertDetails: []report.AlertDetail{
			{
				BasicDetail: &report.BasicDetail{
					References:  refs,
					CreatedTime: oci.Created.Unix(),
					Env:         oci.Config.Env,
					Entrypoint:  oci.Config.Entrypoint,
					Cmd:         oci.Config.Cmd,
					WorkingDir:  oci.Config.WorkingDir,
					Author:      oci.Author,
				},
			},
		},
	}

	err = report.DefaultReportClient().Report(evt)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCommand.AddCommand(cmd.MapImageCommand(scanCommand, scan))
	rootCommand.AddCommand(cmd.NewInfoCommand(plugin.Manifest{
		Name:        "veinmind-basic",
		Author:      "veinmind-team",
		Description: "veinmind-basic scan image basic info",
	}))
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
