package report

import (
	"io/fs"
	"strconv"
	"syscall"
	"time"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/group"
	"github.com/chaitin/veinmind-common-go/passwd"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-sensitive/rule"
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
		Ctim: int64(sys.Ctim.Sec),
		Mtim: int64(sys.Mtim.Sec),
		Atim: int64(sys.Mtim.Sec),
	}, nil
}

func GenerateSensitiveFileEvent(path string, rule rule.SensitiveRule, info fs.FileInfo, image api.Image) (*report.ReportEvent, error) {
	fDetail, err := file2FileDetail(info, path)
	if err != nil {
		return nil, err
	}

	// parse passwd info
	entries, err := passwd.ParseImagePasswd(image)
	if err != nil {
		log.Error(err)
	} else {
		for _, e := range entries {
			if strconv.FormatInt(fDetail.Uid, 10) == e.Uid {
				fDetail.Uname = e.Username
				break
			}
		}
	}

	// parse group info
	gEntries, err := group.ParseImageGroup(image)
	if err != nil {
		log.Error(err)
	} else {
		for _, ge := range gEntries {
			if strconv.FormatInt(fDetail.Gid, 10) == ge.Gid {
				fDetail.Gname = ge.GroupName
				break
			}
		}
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
