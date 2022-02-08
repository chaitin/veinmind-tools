package scanner_common

import "github.com/chaitin/veinmind-tools/veinmind-malicious/database/model"

type ScanPlugin interface {
	Scan(opt ScanOption) (model.ReportData, error)
	PluginName() string
}
