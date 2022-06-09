package report

import (
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-sensitive/rule"
	"github.com/chaitin/veinmind-tools/veinmind-common/go/service/report"
	"io/fs"
	"syscall"
	"time"
)

func localRuleLevel2EventLevel(level string) report.Level {
	switch level {
	case "critical":
		return report.Critical
	case "high":
		return report.High
	case "medium":
		return report.Medium
	case "low":
		return report.Low
	}

	return report.None
}

func file2FileDetail(info fs.FileInfo, path string) (report.FileDetail, error) {
	sys := info.Sys().(*syscall.Stat_t)

	return report.FileDetail{
		Path: path,
		Perm: info.Mode(),
		Size: info.Size(),
		Uid:  int64(sys.Uid),
		Gid:  int64(sys.Gid),
		Ctim: sys.Ctim.Sec,
		Mtim: sys.Mtim.Sec,
		Atim: sys.Mtim.Sec,
	}, nil
}

func GenerateSensitiveFileEvent(path string, rule rule.SensitiveRule, info fs.FileInfo, image api.Image) (*report.ReportEvent, error) {
	fDetail, err := file2FileDetail(info, path)
	if err != nil {
		return nil, err
	}

	r := &report.ReportEvent{
		ID:         image.ID(),
		Time:       time.Now(),
		Level:      localRuleLevel2EventLevel(rule.Level),
		DetectType: report.Image,
		EventType:  report.Risk,
		AlertType:  report.Sensitive,
		AlertDetails: []report.AlertDetail{
			{
				SensitiveFileDetail: &report.SensitveFileDetail{
					FileDetail:      fDetail,
					RuleID:          rule.Id,
					RuleName:        rule.Name,
					RuleDescription: rule.Description,
				},
			},
		},
	}

	return r, nil
}
