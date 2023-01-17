package analyzer

import (
	"context"
	"encoding/json"
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
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-vuln/model"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-vuln/sdk/osv"
	"golang.org/x/sync/semaphore"
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

func TransferAsset(res model.ScanResult) *report.AssetDetail {
	assetDetail := &report.AssetDetail{
		OS: report.AssetOSDetail{
			Family: res.OSInfo.Family,
			Name:   res.OSInfo.Name,
			Eosl:   res.OSInfo.Eosl,
		},
		PackageInfos: transferPackage(res),
		Applications: transferApplication(res),
	}
	return assetDetail
}

func TransferVuln(res model.ScanResult) report.GeneralDetail {
	vulnList := make([]osv.Vulnerability, 0)
	for _, pkgInfo := range res.PackageInfos {
		for _, pkg := range pkgInfo.Packages {
			vulnList = append(vulnList, pkg.Vulnerabilities...)
		}
	}
	for _, app := range res.Applications {
		for _, pkg := range app.Libraries {
			vulnList = append(vulnList, pkg.Vulnerabilities...)
		}
	}
	data, err := json.Marshal(vulnList)
	if err != nil {
		log.Error(err)
		return nil
	}
	return data
}

func transferPackage(res model.ScanResult) []report.AssetPackageDetails {
	var assetPackageDetailsList []report.AssetPackageDetails
	var assetPackageDetails []report.AssetPackageDetail

	for _, pkgInfo := range res.PackageInfos {
		for i, pkg := range pkgInfo.Packages {
			// temp arch format
			if pkg.Arch == "" || pkg.Arch == "None" {
				pkgInfo.Packages[i].Arch = defaultArch
			}
			assetPackageDetails = append(assetPackageDetails, report.AssetPackageDetail{
				Name:       pkg.Name,
				Version:    pkg.Version,
				Release:    pkg.Release,
				Epoch:      pkg.Epoch,
				Arch:       pkgInfo.Packages[i].Arch,
				SrcName:    pkg.SrcName,
				SrcEpoch:   pkg.SrcEpoch,
				SrcRelease: pkg.SrcRelease,
				SrcVersion: pkg.SrcVersion,
			})
		}
		assetPackageDetailsList = append(assetPackageDetailsList, report.AssetPackageDetails{
			FilePath: pkgInfo.FilePath,
			Packages: assetPackageDetails,
		})
		assetPackageDetails = []report.AssetPackageDetail{}
	}

	return assetPackageDetailsList
}

func transferApplication(res model.ScanResult) []report.AssetApplicationDetails {
	var assetApplicationDetailsList []report.AssetApplicationDetails
	var assetPackageDetails []report.AssetPackageDetail

	for _, app := range res.Applications {
		for i, pkg := range app.Libraries {
			// temp arch format
			if pkg.Arch == "" || pkg.Arch == "None" {
				app.Libraries[i].Arch = defaultArch
			}
			assetPackageDetails = append(assetPackageDetails, report.AssetPackageDetail{
				Name:       pkg.Name,
				Version:    pkg.Version,
				Release:    pkg.Release,
				Epoch:      pkg.Epoch,
				Arch:       app.Libraries[i].Arch,
				SrcName:    pkg.SrcName,
				SrcEpoch:   pkg.SrcEpoch,
				SrcRelease: pkg.SrcRelease,
				SrcVersion: pkg.SrcVersion,
			})
		}
		assetApplicationDetailsList = append(assetApplicationDetailsList, report.AssetApplicationDetails{
			Type:     app.Type,
			FilePath: app.FilePath,
			Packages: assetPackageDetails,
		})
		assetPackageDetails = []report.AssetPackageDetail{}
	}

	return assetApplicationDetailsList
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
