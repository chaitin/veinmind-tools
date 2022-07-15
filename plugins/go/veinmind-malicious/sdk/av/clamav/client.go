package clamav

import (
	"errors"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/sdk/av"
	"github.com/dutchcoders/go-clamd"
	"io"
	"strings"
)

func New(clamavHost, clamavPort string) ClamavAddress {
	return ClamavAddress{ClamavHost: clamavHost, ClamavPort: clamavPort, ClamavConnect: clamd.NewClamd("tcp://" + clamavHost + ":" + clamavPort)}
}

func (self *ClamavAddress) Active() bool {
	if self.ClamavConnect == nil {
		return false
	} else {
		if self.ClamavConnect.Ping() != nil {
			return false
		} else {
			return true
		}
	}
}

func (self *ClamavAddress) ScanStream(stream io.Reader) ([]av.ScanResult, error) {
	abort := make(chan bool, 1)
	response, err := self.ClamavConnect.ScanStream(stream, abort)
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

type ClamavAddress struct {
	ClamavHost    string
	ClamavPort    string
	ClamavConnect *clamd.Clamd
}

type SizeLimitReachedError struct {
}

func (self *SizeLimitReachedError) Error() string {
	return "File Size Limit Reached"
}

type ResultParseError struct {
}

func (self *ResultParseError) Error() string {
	return "Clamav Result Parse Error"
}
