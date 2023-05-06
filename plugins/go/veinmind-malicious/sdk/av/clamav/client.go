package clamav

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/testwill/go-clamd"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/sdk/av"
)

var client = func() *clamd.Clamd {
	var CLAMD_ADDRESS = "tcp://" + os.Getenv("CLAMD_HOST") + ":" + os.Getenv("CLAMD_PORT")
	c := clamd.NewClamd(CLAMD_ADDRESS)
	return c
}()

func Active() bool {
	if client == nil {
		return false
	} else {
		if client.Ping() != nil {
			return false
		} else {
			return true
		}
	}
}

func ScanStream(stream io.Reader) ([]av.ScanResult, error) {
	abort := make(chan bool, 1)
	response, err := client.ScanStream(stream, abort)
	defer func() {
		close(abort)
	}()
	if err != nil {
		if strings.Contains(err.Error(), "broken pipe") {
			return nil, new(SizeLimitReachedError)
		}

		return nil, err
	}
	ret := make([]clamd.ScanResult, 0, len(response))
	for s := range response {
		if s.Status == clamd.RES_FOUND {
			ret = append(ret, *s)
		} else if s.Status == clamd.RES_ERROR {
			return nil, errors.New(s.Description)
		} else if s.Status == clamd.RES_PARSE_ERROR {
			return nil, new(ResultParseError)
		}
	}

	// 转换为公共结构体
	retCommon := []av.ScanResult{}
	for _, r := range ret {
		commonResult := av.ScanResult{
			EngineName:  "ClamAV",
			Description: r.Description,
			IsMalicious: true,
			Method:      "blacklist",
		}

		retCommon = append(retCommon, commonResult)
	}

	return retCommon, nil
}

type ServiceInfo struct {
	Version  string
	Pools    string
	State    string
	Threads  string
	Memstats string
	Queue    string
}

type SizeLimitReachedError struct {
}

func (s *SizeLimitReachedError) Error() string {
	return "File Size Limit Reached"
}

type ResultParseError struct {
}

func (r *ResultParseError) Error() string {
	return "Clamav Result Parse Error"
}
