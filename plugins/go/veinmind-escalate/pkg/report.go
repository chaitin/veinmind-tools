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
	CAPPATTERN        string    = "CapEff:\\s*?[a-z0-9]+\\s"
	KERNELPATTERN     string    = "([0-9]{1,})\\.([0-9]{1,})\\.([0-9]{1,})-[0-9]{1,}-[a-zA-Z]{1,}"
	SUDOREGEX         string    = "(\\w{1,})\\s\\w{1,}=\\(.*\\)\\s(.*)"
	CVEREASON         string    = "Your system has an insecure kernel version that is affected by a CVE vulnerability:"
	DOCKERAPIREASON   string    = "Docker remote API is opened which is can be used for escalating"
	SUDOREASON        string    = "This file is granted sudo privileges and can be used for escalating,you can check it in /etc/sudoers"
	MOUNTREASON       string    = "There are some sensitive files or directory mounted"
	READREASON        string    = "This file is sensitive and is readable to all users"
	WRITEREASON       string    = "This file is sensitive and is writable to all users"
	SUIDREASON        string    = "This file is granted suid privileges and belongs to root. And this file can be interacted with, there is a risk of elevation"
	EMPTYPASSWDREASON string    = "This user is privileged but does not have a password set"
	CAPREASON         string    = "There are unsafe linux capability granted"
)

func AddResult(path string, reason string, detail string) {
	result := &models.EscalateResult{
		Target: path,
		Reason: reason,
		Detail: detail,
	}
	escalateLock.Lock()
	res = append(res, result)
	escalateLock.Unlock()
}

func GenerateContainerRoport(container api.Container) error {
	if len(res) > 0 {
		detail, err := json.Marshal(res)
		if err == nil {
			Reportevent := report.ReportEvent{
				ID:         container.ID(),
				Time:       time.Now(),
				Level:      report.High,
				DetectType: report.Container,
				EventType:  report.Risk,
				AlertType:  report.Escalate,
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
func FileClose(file api.File, err error) {
	if err == nil {
		file.Close()
	}
}
