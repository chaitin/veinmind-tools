package clamav

import (
	"errors"
	"github.com/dutchcoders/go-clamd"
	"io"
	"strings"
)

func ScanFile(address string, path string) ([]clamd.ScanResult, error) {
	c := clamd.NewClamd(address)
	response, err := c.ScanFile(path)
	if err != nil {
		return nil, err
	}
	ret := make([]clamd.ScanResult, 0, len(response))
	for s := range response {
		if s.Status == clamd.RES_FOUND {
			ret = append(ret, *s)
		} else if s.Status == clamd.RES_ERROR || s.Status == clamd.RES_PARSE_ERROR {
			return nil, errors.New(s.Description)
		}
	}
	return ret, nil
}

func ScanStream(address string, stream io.Reader) ([]clamd.ScanResult, error) {
	c := clamd.NewClamd(address)
	abort := make(chan bool, 1)
	response, err := c.ScanStream(stream, abort)
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
	return ret, nil
}

type ServiceInfo struct {
	Version  string
	Pools    string
	State    string
	Threads  string
	Memstats string
	Queue    string
}

func QueryServiceInfo(sockPath string) (ServiceInfo, error) {
	c := clamd.NewClamd(sockPath)
	var r ServiceInfo
	response, err := c.Version()
	if err != nil {
		return r, err
	}
	for s := range response {
		r = ServiceInfo{Version: s.Raw}
		break
	}
	stats, err := c.Stats()
	if err != nil {
		return r, err
	}
	r.Pools = stats.Pools
	r.State = stats.State
	r.Threads = stats.Threads
	r.Memstats = stats.Memstats
	r.Queue = stats.Queue
	return r, nil
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
