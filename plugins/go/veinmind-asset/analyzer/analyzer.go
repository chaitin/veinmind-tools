package analyzer

import (
	"context"
	"io/fs"
	"sync"

	dio "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/trivy/pkg/fanal/analyzer"
	_ "github.com/aquasecurity/trivy/pkg/fanal/analyzer/all"
	"github.com/aquasecurity/trivy/pkg/fanal/artifact"
	"github.com/aquasecurity/trivy/pkg/fanal/types"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-asset/model"
	"golang.org/x/sync/semaphore"
)

func ScanImage(image api.Image, parallel int64) (model.ScanImageResult, error) {

	ctx := context.Background()
	limit := semaphore.NewWeighted(parallel)

	// TODO 扫描配置开放
	var artifactOpt artifact.Option
	var analysisOpt analyzer.AnalysisOptions

	ag := analyzer.NewAnalyzerGroup(artifactOpt.AnalyzerGroup, artifactOpt.DisabledAnalyzers)
	var wg sync.WaitGroup
	res := new(analyzer.AnalysisResult)

	image.Walk("/", func(path string, info fs.FileInfo, err error) error {
		open := func() (dio.ReadSeekCloserAt, error) {
			file, err := image.Open(path)
			if err != nil {
				return nil, err
			}
			return file, nil
		}

		ag.AnalyzeFile(ctx, &wg, limit, res, "", path, info, open, nil, analysisOpt)
		return nil
	})
	wg.Wait()
	res.Sort()
	blobInfo := types.BlobInfo{
		SchemaVersion: types.BlobJSONSchemaVersion,
		Digest:        "",
		OS:            res.OS,
		Repository:    res.Repository,
		PackageInfos:  res.PackageInfos,
		Applications:  res.Applications,
		// SystemFiles:     res.SystemInstalledFiles,
		OpaqueDirs:      []string{},
		WhiteoutFiles:   []string{},
		CustomResources: res.CustomResources,

		// For Red Hat
		BuildInfo: res.BuildInfo,
	}

	//result = append(result, blobInfo)
	//artifactDetail := applier.ApplyLayers(result)
	return parseResults(image, blobInfo)
}

func parseResults(image api.Image, res types.BlobInfo) (model.ScanImageResult, error) {
	name, err := image.Repos()
	id := image.ID()
	if err != nil {
		return model.ScanImageResult{}, err
	}
	if len(name) == 0 {
		name = append(name, id)
	}
	os := types.OS{}
	// os判空，否则可能会导致空指针
	if res.OS != nil {
		os = *res.OS
	}
	// format results
	scanRes := model.ScanImageResult{
		ImageName:    name[0],
		ImageID:      id,
		ImageOSInfo:  os,
		PackageTotal: len(res.PackageInfos),
		PackageInfos: res.PackageInfos,
		ApplicationTotal: func() int {
			var a int
			for _, lib := range res.Applications {
				a += len(lib.Libraries)
			}
			return a
		}(),
		Applications: res.Applications,
	}

	return scanRes, nil
}
