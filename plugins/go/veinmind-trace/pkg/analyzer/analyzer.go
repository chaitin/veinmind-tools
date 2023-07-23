package analyzer

import (
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report/event"
)

var Group = make([]Analyzer, 0)

type Analyzer interface {
	Scan(container api.Container)
	Result() []*event.TraceEvent
}
