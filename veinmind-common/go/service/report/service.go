package report

import (
	"context"
	"github.com/chaitin/libveinmind/go/plugin/service"
	"golang.org/x/sync/errgroup"
)

const Namespace = "github.com/chaitin/veinmind-tools/veinmind-common/go/report"
const BufferSize = 1 << 8


type ReportService struct {
	EventChannel chan ReportEvent
}

type reportClient struct {
	ctx    context.Context
	group  *errgroup.Group
	Report func([]ReportEvent) error
}

func (s *ReportService) Report(evts []ReportEvent){
	for _, evt := range evts {
		s.EventChannel <- evt
	}
}

func (s *ReportService) Add(registry *service.Registry) {
	registry.Define(Namespace, struct{}{})
	registry.AddService(Namespace, "report", s.Report)
}

func NewReportService() *ReportService {
	return &ReportService{
		EventChannel: make(chan ReportEvent, BufferSize),
	}
}
