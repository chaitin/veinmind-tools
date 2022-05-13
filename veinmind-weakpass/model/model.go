package model

type WeakpassType int

const (
	SSH WeakpassType = iota
	TOMCAT
	REDIS
)

func (self *WeakpassType) ToString() string {
	switch *self {
	case SSH:
		return "SSH"
	case TOMCAT:
		return "TOMCAT"
	case REDIS:
		return "REDIS"
	}

	return ""
}

type ScanImageResult struct {
	// 镜像名称
	ImageName string

	// 镜像ID
	ImageID string

	// 弱口令结果
	WeakpassResults []WeakpassResult
}

type WeakpassResult struct {
	// 弱口令类型
	PassType WeakpassType

	// 弱口令账户
	Username string

	// 弱口令
	Password string

	// 弱口令位置
	Filepath string
}
