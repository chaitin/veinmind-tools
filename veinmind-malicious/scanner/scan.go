package scanner

import (
	"github.com/chaitin/veinmind-tools/veinmind-malicious/database/model"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/scanner/malicious"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/scanner/scanner_common"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/sdk/common"
)

var ScanPlugins = func() map[string]scanner_common.ScanPlugin {
	s := make(map[string]scanner_common.ScanPlugin)

	// Malicious
	maliciousPlugin := malicious.MaliciousPlugin{}
	s[maliciousPlugin.PluginName()] = &maliciousPlugin

	return s
}()

var ScanPluginsName = func() []string {
	s := []string{}
	for n, _ := range ScanPlugins {
		s = append(s, n)
	}

	return s
}()

//TODO: 逻辑待实现
func MergeReportData(datas []model.ReportData) model.ReportData {
	if len(datas) >= 1 {
		return datas[0]
	}
	return model.ReportData{}
}

func Scan(opt scanner_common.ScanOption) (scanReportAll model.ReportData, err error) {
	datas := []model.ReportData{}
	for _, pluginName := range opt.EnablePlugins {
		if p, ok := ScanPlugins[pluginName]; ok {
			data, err := p.Scan(opt)
			if err != nil {
				common.Log.Error(err)
			}

			datas = append(datas, data)
		}
	}

	return MergeReportData(datas), nil
}
