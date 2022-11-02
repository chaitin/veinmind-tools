package engine

import (
	"encoding/json"
	"path/filepath"
	"time"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report"

	selfreport "github.com/chaitin/veinmind-tools/plugins/go/veinmind-unsafe-mount/pkg/report"
)

func DetectContainerUnsafeMount(container api.Container) (events []report.ReportEvent, err error) {
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
				eBytes, err := json.Marshal(selfreport.Event{
					Source:      mount.Source,
					Destination: mount.Destination,
					Type:        mount.Type,
				})
				if err != nil {
					continue
				}

				events = append(events, report.ReportEvent{
					ID:             container.ID(),
					Time:           time.Now(),
					Level:          report.High,
					DetectType:     report.Container,
					EventType:      report.Risk,
					GeneralDetails: []report.GeneralDetail{eBytes},
				})
			}
		}
	}

	return
}
