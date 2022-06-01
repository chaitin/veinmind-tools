package model

type ScanImageResult struct {
	// 镜像名称
	ImageName string

	// 镜像ID
	ImageID string

	// 服务名称
	ServiceName string

	// 弱口令结果
	WeakpassResults []WeakpassResult
}

type WeakpassResult struct {
	Username string
	Password string
	Filepath string
}

type Config struct {
	Username string
	Dictpath string
	Thread   int
}
