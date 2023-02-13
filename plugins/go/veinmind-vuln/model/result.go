package model

import (
	"github.com/aquasecurity/trivy/pkg/fanal/types"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-vuln/sdk/osv"
)

type ScanResult struct {
	// 镜像/容器名称
	Name string

	// 镜像/容器ID
	ID string

	// 系统信息
	OSInfo *types.OS

	// 系统pkg总数
	PackageTotal int

	// 系统pkg详情
	PackageInfos []PackageInfo

	// 应用依赖总数
	ApplicationTotal int

	// 应用依赖详情
	Applications []Application

	// 漏洞总数
	CveTotal int
}

type PackageInfo struct {
	FilePath string
	Packages []Package
}

type Application struct {
	Type      string
	FilePath  string `json:",omitempty"`
	Libraries []Package
}

type Package struct {
	types.Package
	Vulnerabilities []osv.Vulnerability
}
