package detect

import (
	"io/fs"
	"syscall"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report"
)

func Convert2ReportEvent(image api.Image, info FileInfo, res Result) (*report.ReportEvent, error) {
	if res.Data.RiskLevel == 0 {
		return nil, nil
	}

	var reportLevel report.Level
	switch level := res.Data.RiskLevel; {
	case level <= 5:
		reportLevel = report.Low
		break
	case level <= 10:
		reportLevel = report.Medium
		break
	case level <= 15:
		reportLevel = report.High
		break
	case level <= 20:
		reportLevel = report.Critical
		break
	default:
		return nil, nil
	}

	fileDetail, err := file2FileDetail(info.RawFileInfo, info.Path)
	if err != nil {
		return nil, err
	}

	return &report.ReportEvent{
		ID:         image.ID(),
		Level:      reportLevel,
		DetectType: report.Image,
		EventType:  report.Invasion,
		AlertType:  report.Backdoor,
		AlertDetails: []report.AlertDetail{
			{
				WebshellDetail: &report.WebshellDetail{
					FileDetail: fileDetail,
					Type:       res.Data.Type,
					Engine:     res.Data.Engine,
					Reason:     res.Data.Reason,
				},
			},
		},
	}, nil
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
