package main

import (
	"io"
	"runtime"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gogf/gf/text/gstr"
	"golang.org/x/sync/errgroup"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-sensitive/cache"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-sensitive/report"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-sensitive/rule"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-sensitive/veinfs"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-sensitive/vregex"
)

var (
	defaultLimit         = 5
	defaultContextLength = 500
)

func init() {
	limit := runtime.NumCPU() * 5
	if limit > defaultLimit {
		defaultLimit = limit
	}
}

func Scan(c *cmd.Command, image api.Image) (err error) {
	conf := rule.SingletonConf()

	eg := errgroup.Group{}
	eg.SetLimit(defaultLimit)

	count := uint64(0)
	log.Infof("%s scan file start", image.ID())
	veinfs.Walk(image, "/", func(info *veinfs.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		count += 1
		eg.Go(func() error {
			return scan(image, info.Path, info, conf)
		})

		return nil
	})
	eg.Wait()

	log.Infof("%s scan file count %d", image.ID(), count)
	return nil
}

func scan(image api.Image, path string, info *veinfs.FileInfo, conf *rule.Config) error {
	// check white path cache
	if cache.WhitePath.Contains(path) {
		return nil
	}

	// check path rule cache
	rules, ok := cache.PathRule.Get(path)
	if ok {
		if len(rules) > 0 {
			for _, r := range rules {
				reportEvent(info.Path, r, info, image, "", nil)
			}
		}
	} else {
		for _, pattern := range conf.WhiteList.PathPattern {
			if pattern != "" && vregex.IsMatchString(pattern, info.Path) {
				cache.WhitePath.Add(path)
				return nil
			}
		}

		for _, r := range conf.Rule {
			if r.FilePathPattern != "" && vregex.IsMatchString(r.FilePathPattern, info.Path) {
				cache.PathRule.SetOrAppend(path, r)
				reportEvent(info.Path, r, info, image, "", nil)
			}
		}
	}

	// check file hash cache
	rules, ok = cache.HashRule.Get(info.Sha256)
	if ok {
		for _, r := range rules {
			reportEvent(info.Path, r, info, image, "", nil)
		}
		return nil
	}

	// check file type
	if info.ELF {
		cache.HashRule.Set(info.Sha256, map[int64]rule.Rule{})
		return nil
	}

	// match mime type
	mimeMatch := false
	fp, err := image.Open(path)
	if err != nil {
		log.Error(image.ID(), path, err)
		return nil
	}
	defer fp.Close()

	m, err := mimetype.DetectReader(fp)
	if err != nil {
		log.Warn(image.ID(), path, err)
		return nil
	}

	if gstr.HasPrefix(m.String(), "text/") {
		mimeMatch = true
	} else {
		for mime := range conf.MIMEMap {
			if m.String() == mime {
				mimeMatch = true
			}
		}
	}
	if !mimeMatch {
		cache.HashRule.Set(info.Sha256, make(map[int64]rule.Rule))
		return nil
	}

	_, err = fp.Seek(0, io.SeekStart)
	if err != nil {
		log.Error(image.ID(), path, err)
		return nil
	}

	data, err := io.ReadAll(fp)
	if err != nil {
		log.Warn(image.ID(), path, err)
		return nil
	}

	for _, r := range conf.Rule {
		if r.MatchPattern == "" {
			continue
		}

		content, loc := vregex.FindIndexWithContextContent(r.MatchPattern, data, defaultContextLength)
		if content == nil || loc == nil || len(loc) != 2 {
			continue
		}

		cache.HashRule.SetOrAppend(info.Sha256, r)
		reportEvent(path, r, info, image, string(content), []int64{int64(loc[0]), int64(loc[1])})
	}

	cache.HashRule.Set(info.Sha256, make(map[int64]rule.Rule))
	return nil
}

func reportEvent(path string, r rule.Rule, info *veinfs.FileInfo, image api.Image, contextContent string, contextContentHighlightLocation []int64) {
	evt, err := report.GenerateSensitiveFileEvent(path, r, info, image, contextContent, contextContentHighlightLocation)
	if err != nil {
		log.Error(image.ID(), path, err)
		return
	}

	err = reportService.Client.Report(evt)
	if err != nil {
		log.Error(image.ID(), path, err)
	}
}
