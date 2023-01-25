package reporter

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"sync"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/containerd"
	"github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/libveinmind/go/iac"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/log"
)

type ReportEvent = report.ReportEvent

type Reporter struct {
	ctx          context.Context // ctx control reporter lifecycle
	cancel       context.CancelFunc
	EventChannel chan report.ReportEvent
	events       []ReportEvent
	eventsMutex  sync.RWMutex
	closeOnce    sync.Once
}

func NewReporter(ctx context.Context) (*Reporter, error) {
	ctx, cancel := context.WithCancel(ctx)
	r := &Reporter{
		ctx:          ctx,
		cancel:       cancel,
		EventChannel: make(chan report.ReportEvent, 1<<8),
		events:       []ReportEvent{},
	}
	return r, nil
}

func (r *Reporter) Listen() {
	defer func() {
		log.GetModule(log.ReporterModuleKey).Info("stop reporter listen")
	}()
	for {
		select {
		case evt := <-r.EventChannel:
			r.eventsMutex.Lock()
			r.events = append(r.events, evt)
			r.eventsMutex.Unlock()
		case <-r.ctx.Done():
			return
		}
	}
}

func (r *Reporter) Close() {
	r.closeOnce.Do(func() {
		r.cancel()
		close(r.EventChannel)

		// sync previous sent event
		r.eventsMutex.Lock()
		for e := range r.EventChannel {
			r.events = append(r.events, e)
		}
		r.eventsMutex.Unlock()
	})
}

func (r *Reporter) Write(writer io.Writer) error {
	// read lock
	r.eventsMutex.RLock()
	defer r.eventsMutex.RUnlock()

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

// Render ReportEvent as table format
func (r *Reporter) Render(writer io.Writer) error {
	// read lock
	r.eventsMutex.RLock()
	defer r.eventsMutex.RUnlock()

	// classify cloud-native objects
	category := make(map[report.DetectType][]ReportEvent)
	for _, e := range r.events {
		if category[e.DetectType] == nil {
			category[e.DetectType] = make([]ReportEvent, 0)
		}

		category[e.DetectType] = append(category[e.DetectType], e)
	}

	// set color
	var (
		white  *color.Color
		blue   *color.Color
		yellow *color.Color
		red    *color.Color
	)
	blue = color.New(color.FgBlue)
	blue.EnableColor()
	yellow = color.New(color.FgYellow)
	yellow.EnableColor()
	red = color.New(color.FgRed)
	red.EnableColor()
	white = color.New(color.FgWhite)
	white.EnableColor()

	// color mapping level function
	colorFns := map[report.Level]func(a ...interface{}) string{}
	colorFns[report.Low] = blue.SprintFunc()
	colorFns[report.Medium] = yellow.SprintFunc()
	colorFns[report.High] = red.SprintFunc()
	colorFns[report.Critical] = red.SprintFunc()
	colorFns[report.None] = white.SprintFunc()

	// render image object
	if v, ok := category[report.Image]; ok {
		t := table.NewWriter()
		t.SetOutputMirror(writer)
		t.AppendHeader(table.Row{"Image", "Severity", "Detail"})
		for _, e := range v {
			// handle object
			obj, err := e.RelatedObject()
			if err != nil {
				log.GetModule(log.ReporterModuleKey).Error(err)
				continue
			}

			image, ok := obj.(api.Image)
			if !ok {
				log.GetModule(log.ReporterModuleKey).Error(errors.New("report: can't cast to image object"))
				continue
			}

			var ref string
			refs, _ := image.RepoRefs()
			if len(refs) > 0 {
				ref = refs[0]
			} else {
				ref = image.ID()
			}

			for _, detail := range e.AlertDetails {
				marshalled, err := json.MarshalIndent(detail, "", "	")
				if err != nil {
					log.GetModule(log.ReporterModuleKey).Error(err)
					continue
				}

				colorFn, ok := colorFns[e.Level]
				if !ok {
					log.GetModule(log.ReporterModuleKey).Error(err)
					continue
				}
				t.AppendRow(table.Row{ref, colorFn(e.Level.String()), string(marshalled)})
			}
		}
		t.Render()
	}

	// render container object
	if v, ok := category[report.Container]; ok {
		t := table.NewWriter()
		t.SetOutputMirror(writer)
		t.AppendHeader(table.Row{"Container", "Image", "Severity", "Detail"})
		for _, e := range v {
			// handle object
			obj, err := e.RelatedObject()
			if err != nil {
				log.GetModule(log.ReporterModuleKey).Error(err)
				continue
			}

			container, ok := obj.(api.Container)
			if !ok {
				log.GetModule(log.ReporterModuleKey).Error(errors.New("report: can't cast to container object"))
				continue
			}

			// relate container image
			var (
				image          api.Image
				imageReference string
			)
			switch cast := container.(type) {
			case *docker.Container:
				image, err = cast.Runtime().OpenImageByID(cast.ImageID())
				if err != nil {
					imageReference = "<none>"
				}
			case *containerd.Container:
				image, err = cast.Runtime().OpenImageByID(cast.ImageID())
				if err != nil {
					imageReference = "<none>"
				}
			}
			if image != nil {
				refs, _ := image.RepoRefs()
				if len(refs) > 0 {
					imageReference = refs[0]
				} else {
					imageReference = image.ID()
				}
			}

			for _, detail := range e.AlertDetails {
				marshalled, err := json.MarshalIndent(detail, "", "	")
				if err != nil {
					log.GetModule(log.ReporterModuleKey).Error(err)
					continue
				}

				colorFn, ok := colorFns[e.Level]
				if !ok {
					log.GetModule(log.ReporterModuleKey).Error(err)
					continue
				}
				t.AppendRow(table.Row{container.Name(), imageReference, colorFn(e.Level.String()), string(marshalled)})
			}
		}
		t.Render()
	}

	// render iac object
	if v, ok := category[report.IaC]; ok {
		t := table.NewWriter()
		t.SetOutputMirror(writer)
		t.AppendHeader(table.Row{"Iac", "Type ", "Severity", "Detail"})
		for _, e := range v {
			// handle object
			obj, err := e.RelatedObject()
			if err != nil {
				log.GetModule(log.ReporterModuleKey).Error(err)
				continue
			}

			iac, ok := obj.(iac.IAC)
			if !ok {
				log.GetModule(log.ReporterModuleKey).Error(errors.New("report: can't cast to iac object"))
				continue
			}

			for _, detail := range e.AlertDetails {
				marshalled, err := json.MarshalIndent(detail, "", "	")
				if err != nil {
					log.GetModule(log.ReporterModuleKey).Error(err)
					continue
				}

				colorFn, ok := colorFns[e.Level]
				if !ok {
					log.GetModule(log.ReporterModuleKey).Error(err)
					continue
				}
				t.AppendRow(table.Row{iac.Path, iac.Type.String(), colorFn(e.Level.String()), string(marshalled)})
			}
		}
		t.Render()
	}

	return nil
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

func (r *Reporter) Events() ([]ReportEvent, error) {
	return r.events, nil
}
