package scanner_common

type ScanEngineType int

const (
	Dockerd ScanEngineType = iota
	Containerd
)

type ScanOption struct {
	EngineType    ScanEngineType
	ImageName     string
	EnablePlugins []string
}
