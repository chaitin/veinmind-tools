package scan

import "github.com/chaitin/veinmind-tools/veinmind-weakpass/model"

type EngineType int

const (
	Dockerd EngineType = iota
	Containerd
)

var EngineTypeMap = map[string]EngineType{
	"dockerd":    Dockerd,
	"containerd": Containerd,
}

type ScanOption struct {
	EngineType  EngineType
	ImageName   string
	ScanThreads int
	Username    string
	Dictpath    string
}

type ScanPlugin interface {
	Scan(opt ScanOption) ([]model.ScanImageResult, error)
}
