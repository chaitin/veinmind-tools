package model

type ScanImageResult struct {
	// 镜像名称
	ImageName string

	// 镜像ID
	ImageID string

	// 弱口令类型
	PassType PasswordType

	// 弱口令结果
	WeakpassResults []WeakpassResult
}

type WeakpassResult struct {

	// 弱口令账户
	Username string

	// 弱口令
	Password string

	// 弱口令位置
	Filepath string
}
