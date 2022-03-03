package virustotal

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/VirusTotal/vt-go"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/sdk/av"
	"os"
)


type AnalysisResult struct {
	Category      string `json:"category"`
	EngineName    string `json:"engine_name"`
	EngineUpdate  string `json:"engine_update"`
	EngineVersion string `json:"engine_version"`
	Method        string `json:"method"`
	Result        string `json:"result"`
}

var client = func() *vt.Client{
	apiKey := os.Getenv("VT_API_KEY")
	if apiKey == "" {
		return nil
	}else{
		c := vt.NewClient(apiKey)
		return c
	}
}()

func Active() bool{
	if client == nil {
		return false
	}else{
		return true
	}
}

func ScanSHA256(ctx context.Context, sha256 string)([]av.ScanResult, error){
	retCommon := []av.ScanResult{}
	done := make(chan struct{})

	if client == nil {
		return nil, errors.New("Virustotal Client Init Failed")
	}else{
		// 获取分析结果
		go func() {
			vtFile, err := client.GetObject(vt.URL("files/%s", sha256))
			if err != nil {
				return
			}

			r, err := vtFile.Get("last_analysis_results")
			if err != nil {
				return
			}
			if r == nil {
				return
			}

			// 解析结果
			rMap := r.(map[string]interface{})
			for _, detail := range rMap {
				detailJson, err := json.Marshal(detail)
				if err != nil {
					continue
				}

				analysisResult := AnalysisResult{}
				err = json.Unmarshal(detailJson, &analysisResult)
				if err != nil {
					continue
				}

				if analysisResult.Category == "malicious" {
					commonResult := av.ScanResult{
						Description: analysisResult.Result,
						Method: analysisResult.Method,
						EngineName: analysisResult.EngineName,
						IsMalicious: true,
					}

					retCommon = append(retCommon, commonResult)
				}
			}

			done <- struct{}{}
		}()
	}

	select {
	case <- ctx.Done():
		return retCommon, nil
	case <- done:
		return retCommon, nil
	}
}
