package detect

import (
	"errors"
	"io/fs"
	"syscall"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report/event"
)

func Convert2ReportEvent(fs api.FileSystem, info FileInfo, res Result) (*event.Event, error) {
	if res.Data.RiskLevel == 0 {
		return nil, nil
	}

	var reportLevel event.Level
	switch level := res.Data.RiskLevel; {
	case level <= 5:
		reportLevel = event.Low
		break
	case level <= 10:
		reportLevel = event.Medium
		break
	case level <= 15:
		reportLevel = event.High
		break
	case level <= 20:
		reportLevel = event.Critical
		break
	default:
		return nil, nil
	}

	fileDetail, err := file2FileDetail(info.RawFileInfo, info.Path)
	if err != nil {
		return nil, err
	}

	switch obj := fs.(type) {
	case api.Image:
		return &event.Event{
			BasicInfo: &event.BasicInfo{
				ID:         obj.ID(),
				Object:     event.NewObject(obj),
				Level:      reportLevel,
				DetectType: event.Image,
				EventType:  event.Invasion,
				AlertType:  event.Webshell,
			},
			DetailInfo: &event.DetailInfo{
				AlertDetail: &event.WebshellDetail{
					FileDetail: fileDetail,
					Type:       res.Data.Type,
					Engine:     res.Data.Engine,
					Reason:     res.Data.Reason,
				},
			},
		}, nil
	case api.Container:
		return &event.Event{
			BasicInfo: &event.BasicInfo{
				ID:         obj.ID(),
				Object:     event.NewObject(obj),
				Level:      reportLevel,
				DetectType: event.Container,
				EventType:  event.Invasion,
				AlertType:  event.Webshell,
			},
			DetailInfo: &event.DetailInfo{
				AlertDetail: &event.WebshellDetail{
					FileDetail: fileDetail,
					Type:       res.Data.Type,
					Engine:     res.Data.Engine,
					Reason:     res.Data.Reason,
				},
			},
		}, nil
	}

	return nil, errors.New("not supported")
}

func file2FileDetail(info fs.FileInfo, path string) (event.FileDetail, error) {
	sys := info.Sys().(*syscall.Stat_t)

	return event.FileDetail{
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
