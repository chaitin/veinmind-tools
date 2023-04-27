package engine

import (
	"path/filepath"
	"time"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report/event"
)

const DetectType = "UnsafeMount"

func DetectContainerUnsafeMount(container api.Container) (events []event.Event, err error) {
	spec, err := container.OCISpec()
	if err != nil {
		return nil, err
	}

	for _, mount := range spec.Mounts {
		for _, pattern := range UnsafeMountPaths {
			matched, err := filepath.Match(pattern, mount.Source)
			if err != nil {
				continue
			}

			if matched {
				events = append(events, event.Event{
					BasicInfo: &event.BasicInfo{
						ID:         container.ID(),
						Object:     event.NewObject(container),
						Source:     "veinmind-unsafe-mount",
						Time:       time.Now(),
						Level:      event.High,
						DetectType: event.Container,
						EventType:  event.Risk,
						AlertType:  DetectType,
					},
					DetailInfo: &event.DetailInfo{
						AlertDetail: &event.UnSafeMountDetail{
							Mount: event.MountEvent{
								Source:      mount.Source,
								Destination: mount.Destination,
								Type:        mount.Type,
							},
						},
					},
				})
			}
		}
	}
	return
}
