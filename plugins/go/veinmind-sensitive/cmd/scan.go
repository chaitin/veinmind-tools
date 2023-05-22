package main

import (
	"io"
	"runtime"
	"strings"

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

	// scan env
	log.Infof("%s scan env start", image.ID())
	err = scanEnv(image, conf)
	if err != nil {
		return err
	}

	// scan history
	log.Infof("%s scan docker history start", image.ID())
	err = scanDockerHistory(image, conf)
	if err != nil {
		return err
	}

	// scan filesystem
	log.Infof("%s scan file start", image.ID())
	veinfs.Walk(image, "/", func(info *veinfs.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if (!info.Type.IsRegular() && !info.Type.IsDir()) || info.Size > 3*1024*1024 {
			return nil
		}
		count += 1
		eg.Go(func() error {
			return scanFS(image, info.Path, info, conf)
		})

		return nil
	})
	err = eg.Wait()
	if err != nil {
		return err
	}
	log.Infof("%s scan file count %d", image.ID(), count)
	return nil
}

// scanEnv 扫描环境变量中的敏感信息
func scanEnv(image api.Image, conf *rule.Config) error {
	ocispec, err := image.OCISpecV1()
	if err != nil {
		return err
	}
	for _, env := range ocispec.Config.Env {
		for _, r := range conf.Rule {
			if r.Env != "" && vregex.IsMatchString(r.Env, env) {
				envArr := strings.Split(env, "=")
				if len(envArr) == 2 {
					reportEvent("env", image, r, envArr[0], envArr[1], "", "", nil, "", nil)
				}
			}
		}
	}
	return nil
}

// scanHistory 扫描镜像历史命令中的敏感信息
func scanDockerHistory(image api.Image, conf *rule.Config) error {
	ocispec, err := image.OCISpecV1()
	if err != nil {
		return err
	}

	for _, history := range ocispec.History {
		for _, r := range conf.Rule {
			if r.MatchPattern != "" && vregex.IsMatchString(r.MatchPattern, history.CreatedBy) {
				reportEvent("history", image, r, "", "", history.CreatedBy, "", nil, "", nil)
			}
		}
	}
	return nil
}

// scanFS 扫描镜像文件系统中的敏感信息
func scanFS(image api.Image, path string, info *veinfs.FileInfo, conf *rule.Config) error {
	// check white path cache
	if cache.WhitePath.Contains(path) {
		return nil
	}

	// check path rule cache
	rules, ok := cache.PathRule.Get(path)
	if ok {
		if len(rules) > 0 {
			for _, r := range rules {
				reportEvent("file", image, r, "", "", "", info.Path, info, "", nil)
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
				reportEvent("file", image, r, "", "", "", info.Path, info, "", nil)
			}
		}
	}

	// check file hash cache
	rules, ok = cache.HashRule.Get(info.Sha256)
	if ok {
		for _, r := range rules {
			reportEvent("file", image, r, "", "", "", info.Path, info, "", nil)
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
		if !info.Type.IsDir() {
			log.Error(image.ID(), path, err)
		}
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
		reportEvent("file", image, r, "", "", "", path, info, string(content), []int64{int64(loc[0]), int64(loc[1])})
	}

	cache.HashRule.Set(info.Sha256, make(map[int64]rule.Rule))
	return nil
}

func reportEvent(eventType string, image api.Image, r rule.Rule, envKey string, envValue string, history string, path string, info *veinfs.FileInfo, contextContent string, contextContentHighlightLocation []int64) {
	switch eventType {
	case "env":
		evt, err := report.GenerateSensitiveEnvEvent(image, r, envKey, envValue)
		if err != nil {
			log.Error(image.ID(), path, err)
			return
		}
		err = reportService.Client.Report(evt)
		if err != nil {
			log.Error(image.ID(), path, err)
		}
	case "history":
		evt, err := report.GenerateSensitiveDockerHistoryEvent(image, r, history)
		if err != nil {
			log.Error(image.ID(), path, err)
			return
		}
		err = reportService.Client.Report(evt)
		if err != nil {
			log.Error(image.ID(), path, err)
		}
	case "file":
		evt, err := report.GenerateSensitiveFileEvent(image, r, path, info, contextContent, contextContentHighlightLocation)
		if err != nil {
			log.Error(image.ID(), path, err)
			return
		}
		err = reportService.Client.Report(evt)
		if err != nil {
			log.Error(image.ID(), path, err)
		}
	}
}
