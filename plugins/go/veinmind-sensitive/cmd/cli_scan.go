package main

import (
	"bytes"
	"io"
	"io/fs"
	"io/ioutil"
	"runtime"
	"strings"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin/log"
	reportService "github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-sensitive/report"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-sensitive/rule"
	"github.com/gabriel-vasile/mimetype"
	"golang.org/x/sync/errgroup"
)

func Scan(c *cmd.Command, image api.Image) (err error) {
	conf := rule.SingletonConf()

	eg := errgroup.Group{}
	limit := runtime.NumCPU() * 10
	if limit < 30 {
		limit = 30
	}
	eg.SetLimit(limit)

	num := 0
	image.Walk("/", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Warn(image.ID(), path, err)
			return nil
		}
		if info.IsDir() {
			return nil
		}

		num += 1
		eg.Go(func() error {
			return scan(image, path, info, conf)
		})
		return nil
	})

	log.Info(image.ID(), "scan file count ", num)
	return eg.Wait()
}

func scan(image api.Image, path string, info fs.FileInfo, conf *rule.SensitiveConfig) error {
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
				log.Warn(image.ID(), path, err)
			} else {
				err = reportService.DefaultReportClient().Report(*evt)
				if err != nil {
					log.Error(image.ID(), path, err)
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
		log.Error(image.ID(), path, err)
		return nil
	}
	defer f.Close()

	m, err := mimetype.DetectReader(f)
	if err != nil {
		log.Warn(image.ID(), path, err)
	} else {
		if strings.HasPrefix(m.String(), "text/") {
			mimeMatch = true
		} else {
			for mime := range conf.MIMEMap {
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
			log.Error(image.ID(), path, err)
			return nil
		}

		fb, err = ioutil.ReadAll(f)
		if err != nil {
			log.Warn(image.ID(), path, err)
		}
	} else {
		return nil
	}

	for _, r := range conf.Rules {
		// match content
		if r.MatchRegex != nil && r.MatchRegex.Match(fb) {
			evt, err := report.GenerateSensitiveFileEvent(path, r, info, image)
			if err != nil {
				log.Warn(image.ID(), path, err)
			} else {
				err = reportService.DefaultReportClient().Report(*evt)
				if err != nil {
					log.Error(image.ID(), path, err)
				}
				return nil
			}
		}

		if r.MatchContains != "" && bytes.Contains(fb, []byte(r.MatchContains)) {
			evt, err := report.GenerateSensitiveFileEvent(path, r, info, image)
			if err != nil {
				log.Warn(image.ID(), path, err)
			} else {
				err = reportService.DefaultReportClient().Report(*evt)
				if err != nil {
					log.Error(image.ID(), path, err)
				}
				return nil
			}
		}
	}

	return nil
}
