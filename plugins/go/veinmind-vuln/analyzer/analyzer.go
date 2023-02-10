package analyzer

import (
	"context"
	"io/fs"
	"os"
	"strings"
	"sync"

	dio "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/trivy/pkg/fanal/analyzer"
	_ "github.com/aquasecurity/trivy/pkg/fanal/analyzer/all"
	"github.com/aquasecurity/trivy/pkg/fanal/artifact"
	"github.com/aquasecurity/trivy/pkg/fanal/types"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"golang.org/x/sync/semaphore"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-vuln/model"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-vuln/sdk/osv"
)

var (
	defaultArch = "noarch"
)

func ScanImage(image api.Image, parallel int64) (model.ScanResult, error) {

	res := ScanAsset(image, parallel)
	scanRes := parseResult(res)

	// format results
	name, err := image.RepoRefs()
	id := image.ID()
	if err != nil {
		return model.ScanResult{}, err
	}
	if len(name) == 0 {
		name = append(name, id)
	}
	scanRes.Name = strings.Join(name, ",")
	scanRes.ID = id
	return scanRes, nil
}

func ScanContainer(container api.Container, parallel int64) (model.ScanResult, error) {
	res := ScanAsset(container, parallel)
	scanRes := parseResult(res)
	// format results
	scanRes.Name = container.Name()
	scanRes.ID = container.ID()
	return scanRes, nil
}

func ScanAsset(fileSystem api.FileSystem, parallel int64) *analyzer.AnalysisResult {
	var wg sync.WaitGroup
	ctx := context.Background()
	limit := semaphore.NewWeighted(parallel)
	res := analyzer.NewAnalysisResult()
	opts := analyzer.AnalysisOptions{Offline: true}

	// TODO 扫描配置开放
	var artifactOpt artifact.Option

	ag := analyzer.NewAnalyzerGroup(artifactOpt.AnalyzerGroup, artifactOpt.DisabledAnalyzers)

	fileSystem.Walk("/", func(path string, info fs.FileInfo, err error) error {
		// 如果出现error会导致后续空指针，这里需要处理一下
		if err != nil {
			log.Debug(err)
			return nil
		}

		// Copy From veinmind-malicious, 提速
		// 判断文件类型，跳过特定类型文件
		if (info.Mode() & (os.ModeDevice | os.ModeNamedPipe | os.ModeSocket | os.ModeCharDevice | os.ModeDir)) != 0 {
			log.Debug("Skip: ", path)
			return nil
		}

		open := func() (dio.ReadSeekCloserAt, error) {
			file, err := fileSystem.Open(path)
			if err != nil {
				return nil, err
			}
			return file, nil
		}

		ag.AnalyzeFile(ctx, &wg, limit, res, "", path, info, open, nil, opts)
		return nil
	})

	wg.Wait()
	res.Sort()
	return res
}

func parseResult(scanRes *analyzer.AnalysisResult) model.ScanResult {
	osInfo := &types.OS{}
	if scanRes.OS != nil {
		osInfo = scanRes.OS
	}
	return model.ScanResult{
		OSInfo: osInfo,
		PackageTotal: func() int {
			var a int
			for _, lib := range scanRes.PackageInfos {
				a += len(lib.Packages)
			}
			return a
		}(),
		PackageInfos: parsePackage(scanRes),
		ApplicationTotal: func() int {
			var a int
			for _, lib := range scanRes.Applications {
				a += len(lib.Libraries)
			}
			return a
		}(),
		Applications: parseApplication(scanRes),
	}
}

func parseApplication(scanRes *analyzer.AnalysisResult) []model.Application {
	var appInfos = make([]model.Application, 0)
	for _, apps := range scanRes.Applications {
		info := model.Application{
			FilePath:  apps.FilePath,
			Type:      apps.Type,
			Libraries: make([]model.Package, 0),
		}
		for _, app := range apps.Libraries {
			info.Libraries = append(info.Libraries, model.Package{
				Package:         app,
				Vulnerabilities: make([]osv.Vulnerability, 0),
			})
		}
		if len(info.Libraries) > 0 {
			appInfos = append(appInfos, info)
		}
	}
	return appInfos
}

func parsePackage(scanRes *analyzer.AnalysisResult) []model.PackageInfo {
	var pkgInfos = make([]model.PackageInfo, 0)
	for _, pkgs := range scanRes.PackageInfos {
		info := model.PackageInfo{
			FilePath: pkgs.FilePath,
			Packages: make([]model.Package, 0),
		}
		for _, pkg := range pkgs.Packages {
			info.Packages = append(info.Packages, model.Package{
				Package:         pkg,
				Vulnerabilities: make([]osv.Vulnerability, 0),
			})
		}
		if len(info.Packages) > 0 {
			pkgInfos = append(pkgInfos, info)
		}
	}
	return pkgInfos
}
