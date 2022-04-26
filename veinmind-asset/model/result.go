package model

import "github.com/aquasecurity/fanal/types"

type ScanImageResult struct {
	// 镜像名称
	ImageName string

	// 镜像ID
	ImageID string

	// 系统信息
	ImageInfo types.OS

	// 系统pkg总数
	PackageTotal int

	// 系统pkg详情
	Packages []types.Package

	// 应用依赖总数
	ApplicationTotal int

	// 应用依赖详情
	Applications []types.Application
}
