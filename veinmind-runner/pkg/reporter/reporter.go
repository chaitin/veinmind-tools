package reporter

import (
	"encoding/json"
	"io"

	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report"
)

type ReportEvent struct {
	report.ReportEvent
	ImageRefs []string `json:"image_refs"`
}

type Reporter struct {
	EventChannel chan report.ReportEvent
	closeCh      chan struct{}
	events       []ReportEvent
}

func NewReporter() (*Reporter, error) {
	return &Reporter{
		EventChannel: make(chan report.ReportEvent, 1<<8),
		closeCh:      make(chan struct{}),
		events:       []ReportEvent{},
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

func WriteEvents2Log(events []ReportEvent, writer io.Writer) error {
	if len(events) == 0 {
		return nil
	}

	eventsBytes, err := json.MarshalIndent(events, "", "  ")
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

func (r *Reporter) GetEvents() ([]ReportEvent, error) {
	return r.events, nil
}

func (r *Reporter) convert(event report.ReportEvent) (ReportEvent, error) {
	if event.DetectType == report.IaC {
		return ReportEvent{
			ReportEvent: event,
		}, nil
	}
	// todo: rollback reference or redesigned ReportEvent structure
	return ReportEvent{
		ReportEvent: event,
	}, nil
}
