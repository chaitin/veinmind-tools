package common

type SDK interface {
	GetSDKInfo() (SDKInfo, error)
}

type SDKInfo struct {
}
