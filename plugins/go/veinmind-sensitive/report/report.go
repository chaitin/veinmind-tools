package report

import (
	"strconv"
	"time"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/group"
	"github.com/chaitin/veinmind-common-go/passwd"
	"github.com/chaitin/veinmind-common-go/service/report/event"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-sensitive/rule"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-sensitive/veinfs"
)

func localRuleLevel2EventLevel(level string) event.Level {
	switch level {
	case "critical":
		return event.Critical
	case "high":
		return event.High
	case "medium":
		return event.Medium
	case "low":
		return event.Low
	}

	return event.None
}

func file2FileDetail(info *veinfs.FileInfo, path string) (event.FileDetail, error) {
	return event.FileDetail{
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

func GenerateSensitiveEnvEvent(image api.Image, rule rule.Rule, envKey string, envValue string) (*event.Event, error) {
	r := &event.Event{
		BasicInfo: &event.BasicInfo{
			ID:         image.ID(),
			Object:     event.NewObject(image),
			Source:     "veinmind-sensitive",
			Time:       time.Now(),
			Level:      localRuleLevel2EventLevel(rule.Level),
			DetectType: event.Image,
			EventType:  event.Risk,
			AlertType:  event.SensitiveFile,
		},
		DetailInfo: &event.DetailInfo{
			AlertDetail: &event.SensitiveEnvDetail{
				SensitiveDetail: event.SensitiveDetail{},
				Key:             envKey,
				Value:           envValue,
				RuleID:          rule.Id,
				RuleName:        rule.Name,
				RuleDescription: rule.Description,
			},
		},
	}
	return r, nil
}

func GenerateSensitiveDockerHistoryEvent(image api.Image, rule rule.Rule, history string) (*event.Event, error) {
	r := &event.Event{
		BasicInfo: &event.BasicInfo{
			ID:         image.ID(),
			Object:     event.NewObject(image),
			Source:     "veinmind-sensitive",
			Time:       time.Now(),
			Level:      localRuleLevel2EventLevel(rule.Level),
			DetectType: event.Image,
			EventType:  event.Risk,
			AlertType:  event.SensitiveFile,
		},
		DetailInfo: &event.DetailInfo{
			AlertDetail: &event.SensitiveDockerHistoryDetail{
				SensitiveDetail: event.SensitiveDetail{},
				Value:           history,
				RuleID:          rule.Id,
				RuleName:        rule.Name,
				RuleDescription: rule.Description,
			},
		},
	}
	return r, nil
}

func GenerateSensitiveFileEvent(image api.Image, rule rule.Rule, path string, info *veinfs.FileInfo, contextContent string, contextContentHighlightLocation []int64) (*event.Event, error) {
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

	r := &event.Event{
		BasicInfo: &event.BasicInfo{
			ID:         image.ID(),
			Object:     event.NewObject(image),
			Source:     "veinmind-sensitive",
			Time:       time.Now(),
			Level:      localRuleLevel2EventLevel(rule.Level),
			DetectType: event.Image,
			EventType:  event.Risk,
			AlertType:  event.SensitiveFile,
		},
		DetailInfo: &event.DetailInfo{
			AlertDetail: &event.SensitiveFileDetail{
				FileDetail:                   fDetail,
				RuleID:                       rule.Id,
				RuleName:                     rule.Name,
				RuleDescription:              rule.Description,
				ContextContent:               contextContent,
				ContextContentHighlightRange: contextContentHighlightLocation,
			},
		},
	}

	return r, nil
}
