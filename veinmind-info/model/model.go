package model

import (
	"time"
)

type ImageInfo struct {
	// 镜像ID
	ID string

	// 创建时间
	Created *time.Time

	// 镜像大小
	Size string

	// 镜像版本
	Tag string

	// 环境变量
	Env []string

	// 启动用户
	User string

	// 工作目录
	WorkingDir string

	// CMD
	Cmd []string

	// 可挂载目录
	Volumes map[string]struct{}

	// 入口点
	Entrypoint []string

	// 暴露端口
	ExposedPorts map[string]struct{}

	// 用户信息
	Users []ImageUserInfo
}

type ImageUserInfo struct {
	// 用户名
	Username string

	// UID
	Uid string

	// GID
	Gid string

	// Shell
	Shell string

	// 用户说明
	Description string
}
