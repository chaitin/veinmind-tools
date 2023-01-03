package report

import (
	"strconv"
	"time"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/group"
	"github.com/chaitin/veinmind-common-go/passwd"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-sensitive/rule"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-sensitive/veinfs"
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

func file2FileDetail(info *veinfs.FileInfo, path string) (report.FileDetail, error) {
	return report.FileDetail{
		Path: path,
		Perm: info.Perm,
		Size: int64(info.Size),
		Uid:  int64(info.Uid),
		Gid:  int64(info.Gid),
		Ctim: info.CreateTime.Unix(),
		Mtim: info.ModifyTime.Unix(),
		Atim: info.AccessTime.Unix(),
	}, nil
}

func GenerateSensitiveFileEvent(path string, rule rule.Rule, info *veinfs.FileInfo, image api.Image, contextContent string, contextContentHighlightLocation []int64) (*report.ReportEvent, error) {
	fDetail, err := file2FileDetail(info, path)
	if err != nil {
		return nil, err
	}

	// parse passwd info
	entries, err := passwd.ParseFilesystemPasswd(image)
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
	gEntries, err := group.ParseFilesystemGroup(image)
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
					FileDetail:                   fDetail,
					RuleID:                       rule.Id,
					RuleName:                     rule.Name,
					RuleDescription:              rule.Description,
					ContextContent:               contextContent,
					ContextContentHighlightRange: contextContentHighlightLocation,
				},
			},
		},
	}

	return r, nil
}
