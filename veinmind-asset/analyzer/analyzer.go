package analyzer

import (
	"context"
	"io"
	"io/fs"
	"sync"

	"github.com/aquasecurity/fanal/analyzer"
	_ "github.com/aquasecurity/fanal/analyzer/all"
	"github.com/aquasecurity/fanal/applier"
	"github.com/aquasecurity/fanal/artifact"
	_ "github.com/aquasecurity/fanal/hook/all"
	"github.com/aquasecurity/fanal/types"
	dio "github.com/aquasecurity/go-dep-parser/pkg/io"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/containerd"
	"github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/veinmind-tools/veinmind-asset/model"
	"golang.org/x/sync/semaphore"
)

// TODO 目前api.File 缺失seek方法，导致某些fanal内存在的seek函数无法使用。
// 目前在base_dirs排除了这部分路径
type AtempFile struct {
	api.File
	io.Seeker
}

func ScanImage(image api.Image, parallel int64) (model.ScanImageResult, error) {

	ctx := context.Background()
	limit := semaphore.NewWeighted(parallel)

	// TODO 扫描配置开放
	var artifactOpt artifact.Option
	var analysisOpt analyzer.AnalysisOptions

	// 忽略存在seek的jar扫描和go binary扫描
	// artifactOpt.DisabledAnalyzers = []analyzer.Type{
	disableType := []analyzer.Type{
		analyzer.TypeJar,
		analyzer.TypeGoBinary,
	}

	var result []types.BlobInfo
	ag := analyzer.NewAnalyzerGroup(artifactOpt.AnalyzerGroup, artifactOpt.DisabledAnalyzers)

	switch v := image.(type) {
	case *docker.Image:
		dockerImage := v
		// TODO 增加缓存机制
		for index := 0; index < dockerImage.NumLayers(); index++ {
			layer, err := dockerImage.OpenLayer(index)
			var wg sync.WaitGroup
			res := new(analyzer.AnalysisResult)
			if err == nil {
				layerID := layer.ID()
				// 根目录开始walk
				layer.Walk("/", func(path string, info fs.FileInfo, err error) error {
					open := func() (dio.ReadSeekCloserAt, error) {
						file, err := layer.Open(path)
						if err != nil {
							return nil, err
						}
						return AtempFile{file, nil}, nil
					}

					ag.AnalyzeFile(ctx, &wg, limit, res, "", path, info, open, disableType, analysisOpt)
					return nil
				})
				wg.Wait()
				res.Sort()

				// 将layer扫描结果转化为blobinfo方便Merge
				blobInfo := types.BlobInfo{
					SchemaVersion:   types.BlobJSONSchemaVersion,
					Digest:          "",
					DiffID:          layerID,
					OS:              res.OS,
					Repository:      res.Repository,
					PackageInfos:    res.PackageInfos,
					Applications:    res.Applications,
					SystemFiles:     res.SystemInstalledFiles,
					OpaqueDirs:      []string{},
					WhiteoutFiles:   []string{},
					CustomResources: res.CustomResources,

					// For Red Hat
					BuildInfo: res.BuildInfo,
				}

				result = append(result, blobInfo)
			}
		}
	case *containerd.Image:
		containerdImage := v
		imageID := containerdImage.ID()
		var wg sync.WaitGroup
		res := new(analyzer.AnalysisResult)
		containerdImage.Walk("/", func(path string, info fs.FileInfo, err error) error {
			open := func() (dio.ReadSeekCloserAt, error) {
				file, err := containerdImage.Open(path)
				if err != nil {
					return nil, err
				}
				return AtempFile{file, nil}, nil
			}

			ag.AnalyzeFile(ctx, &wg, limit, res, "", path, info, open, disableType, analysisOpt)
			return nil
		})
		wg.Wait()
		res.Sort()
		blobInfo := types.BlobInfo{
			SchemaVersion:   types.BlobJSONSchemaVersion,
			Digest:          "",
			DiffID:          imageID,
			OS:              res.OS,
			Repository:      res.Repository,
			PackageInfos:    res.PackageInfos,
			Applications:    res.Applications,
			SystemFiles:     res.SystemInstalledFiles,
			OpaqueDirs:      []string{},
			WhiteoutFiles:   []string{},
			CustomResources: res.CustomResources,

			// For Red Hat
			BuildInfo: res.BuildInfo,
		}
		result = append(result, blobInfo)
	}
	artifactDetail := applier.ApplyLayers(result)
	return parseResults(image, artifactDetail)
}

func parseResults(image api.Image, res types.ArtifactDetail) (model.ScanImageResult, error) {
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
		ImageInfo:    os,
		PackageTotal: len(res.Packages),
		Packages:     res.Packages,
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
