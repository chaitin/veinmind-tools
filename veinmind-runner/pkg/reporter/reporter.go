package reporter

import (
	"encoding/json"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/containerd"
	"github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-tools/veinmind-common/go/service/report"
	"github.com/pkg/errors"
	"io"
)

type reportEvent struct {
	report.ReportEvent
	ImageRefs []string `json:"image_refs"`
}

type Reporter struct {
	EventChannel chan report.ReportEvent
	closeCh      chan struct{}
	events       []reportEvent
}

func NewReporter() (*Reporter, error) {
	return &Reporter{
		EventChannel: make(chan report.ReportEvent, 1<<8),
		closeCh:      make(chan struct{}),
		events:       []reportEvent{},
	}, nil
}

func (r *Reporter) Listen() {
	for {
		select {
		case evt := <-r.EventChannel:
			evtN, err := r.convert(evt)
			if err != nil {
				log.Error(err)
			}
			r.events = append(r.events, evtN)
		case <-r.closeCh:
			goto END
		}
	}
END:
	log.Info("Stop reporter listen")
}

func (r *Reporter) StopListen() {
	r.closeCh <- struct{}{}
}

func (r *Reporter) Write(writer io.Writer) error {
	if len(r.events) == 0 {
		return nil
	}

	eventsBytes, err := json.MarshalIndent(r.events, "", "  ")
	if err != nil {
		return err
	}

	_, err = writer.Write(eventsBytes)
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte("\n"))
	return err
}

func (r *Reporter) GetEvents() ([]reportEvent, error) {
	return r.events, nil
}

func (r *Reporter) convert(event report.ReportEvent) (reportEvent, error) {
	dr, _ := docker.New()
	cr, _ := containerd.New()
	runtimes := []api.Runtime{dr, cr}
	var image api.Image
	find := false
	for _, runtime := range runtimes {
		if runtime != nil {
			i, err := runtime.OpenImageByID(event.ID)
			if err != nil {
				continue
			}
			image = i
			find = true
			break
		}
	}
	if !find || image == nil {
		return reportEvent{}, errors.New("Can't get image object")
	}

	refs, err := image.RepoRefs()
	if err != nil {
		refs = []string{}
		log.Error(err)
	}

	//oci, err := image.OCISpecV1()
	//if err != nil {
	//	oci = nil
	//	log.Error(err)
	//}

	return reportEvent{
		ImageRefs:   refs,
		ReportEvent: event,
	}, nil
}
