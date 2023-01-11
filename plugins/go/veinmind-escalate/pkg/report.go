package pkg

import (
	"encoding/json"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-escalate/models"
	"sync"
	"time"
)

var escalateLock sync.Mutex
var res = []*models.EscalateResult{}

const (
	WRITE             checkMode = 2
	READ              checkMode = 4
	MOUNTREASON       string    = "There are some sensitive files or directory mounted"
	CRONFLAG          string    = "INCRON:"
	READREASON        string    = "This file is sensitive and is readable to all users"
	CRONWRITEREASON   string    = "This file appears in the crontab file and is writable to all users"
	WRITEREASON       string    = "This file is sensitive and is writable to all users"
	SUIDREASON        string    = "This file is granted suid privileges and belongs to root. And this file can be interacted with, there is a risk of elevation"
	EMPTYPASSWDREASON string    = "This user is privileged but does not have a password set"
	CAPREASON         string    = "There are unsafe linux capability granted"
)

func AddResult(path string, reason string, detail string) {
	result := &models.EscalateResult{
		Path:   path,
		Reason: reason,
		Detail: detail,
	}
	escalateLock.Lock()
	res = append(res, result)
	escalateLock.Unlock()
}

func GenerateContainerRoport(image api.Container) error {
	if len(res) > 0 {
		detail, err := json.Marshal(res)
		if err == nil {
			Reportevent := report.ReportEvent{
				ID:         image.ID(),
				Time:       time.Now(),
				Level:      report.High,
				DetectType: report.Image,
				EventType:  report.Risk,
				AlertType:  report.Weakpass,
				GeneralDetails: []report.GeneralDetail{
					detail,
				},
			}
			err := report.DefaultReportClient().Report(Reportevent)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
func GenerateImageRoport(image api.Image) error {
	if len(res) > 0 {
		detail, err := json.Marshal(res)
		if err == nil {
			Reportevent := report.ReportEvent{
				ID:         image.ID(),
				Time:       time.Now(),
				Level:      report.High,
				DetectType: report.Image,
				EventType:  report.Risk,
				AlertType:  report.Weakpass,
				GeneralDetails: []report.GeneralDetail{
					detail,
				},
			}
			err := report.DefaultReportClient().Report(Reportevent)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
