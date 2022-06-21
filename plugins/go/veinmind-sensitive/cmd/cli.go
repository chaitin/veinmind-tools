package main

import (
	"bytes"
	"context"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	reportService "github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-sensitive/report"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-sensitive/rule"
	"github.com/gabriel-vasile/mimetype"
	"github.com/spf13/cobra"
	"io"
	"io/fs"
	"io/ioutil"
	"strings"
)

var (
	rootCommand = &cmd.Command{}
	scanCommand = &cmd.Command{
		Use:   "scan image sensitive info",
		Short: "scan image sensitive info",
		PreRun: func(cmd *cobra.Command, args []string) {
			rule.Init()
		},
	}
)

func scan(c *cmd.Command, image api.Image) (err error) {
	conf := rule.SingletonConf()

	image.Walk("/", func(path string, info fs.FileInfo, err error) error {
		// skip white list
		for _, whitePathGlob := range conf.WhiteList.PathsGlob {
			if whitePathGlob.Match(path) {
				return nil
			}
		}

		for _, r := range conf.Rules {
			// match file path
			if r.FilepathRegexp != nil && r.FilepathRegexp.MatchString(path) {
				evt, err := report.GenerateSensitiveFileEvent(path, r, info, image)
				if err != nil {
					log.Error(err)
				} else {
					err = reportService.DefaultReportClient().Report(*evt)
					if err != nil {
						log.Error(err)
					}
					return nil
				}
			}
		}

		// skip not regular file
		if !info.Mode().IsRegular() {
			return nil
		}

		// match mime type
		mimeMatch := false
		f, err := image.Open(path)
		if err != nil {
			log.Error(err)
			return nil
		}
		defer f.Close()

		m, err := mimetype.DetectReader(f)
		if err != nil {
			log.Error(err)
		} else {
			if strings.HasPrefix(m.String(), "text/") {
				mimeMatch = true
			} else {
				for mime, _ := range conf.MIMEMap {
					if m.String() == mime {
						mimeMatch = true
					}
				}
			}
		}

		var fb []byte
		if mimeMatch {
			_, err = f.Seek(0, io.SeekStart)
			if err != nil {
				log.Error(err)
				return nil
			}

			fb, err = ioutil.ReadAll(f)
			if err != nil {
				log.Error(err)
			}
		} else {
			return nil
		}

		for _, r := range conf.Rules {
			// match content
			if r.MatchRegex != nil && r.MatchRegex.Match(fb) {
				evt, err := report.GenerateSensitiveFileEvent(path, r, info, image)
				if err != nil {
					log.Error(err)
				} else {
					err = reportService.DefaultReportClient().Report(*evt)
					if err != nil {
						log.Error(err)
					}
					return nil
				}
			}

			if r.MatchContains != "" && bytes.Contains(fb, []byte(r.MatchContains)) {
				evt, err := report.GenerateSensitiveFileEvent(path, r, info, image)
				if err != nil {
					log.Error(err)
				} else {
					err = reportService.DefaultReportClient().Report(*evt)
					if err != nil {
						log.Error(err)
					}
					return nil
				}
			}
		}

		return nil
	})

	return nil
}

func init() {
	rootCommand.AddCommand(cmd.MapImageCommand(scanCommand, scan))
	rootCommand.AddCommand(cmd.NewInfoCommand(plugin.Manifest{
		Name:        "veinmind-sensitive",
		Author:      "veinmind-team",
		Description: "veinmind-sensitive scan image sensitive data",
	}))
}

func main() {
	if err := rootCommand.ExecuteContext(context.Background()); err != nil {
		panic(err)
	}
}
