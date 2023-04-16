package model

import "github.com/chaitin/veinmind-common-go/service/report/event"

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

// WeakpassResult 弱密码相关信息
type WeakpassResult struct {
	Username    string
	Password    string
	Filepath    string
	ServiceType event.WeakpassService
}

// Record 从文件中解析出来的相关信息
type Record struct {
	Username string
	Password string
	// 除用户名密码外, 有些模块有其他属性
	// 可以记录在此map中
	Attributes map[string]string
}

// Config cli命令中与爆破相关的配置信息
type Config struct {
	Thread   int
	Username string
	Dictpath string
}

// BruteOption tunny 需要的密码爆破相关的信息
// Guess 碰撞的密码
// Records 模块配置文件中提取的密码信息
type BruteOption struct {
	Records Record
	Guess   string
}
