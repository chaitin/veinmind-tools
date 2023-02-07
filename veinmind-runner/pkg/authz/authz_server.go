package authz

import (
	"github.com/chaitin/veinmind-common-go/service/report/event"
	"os"
	"os/signal"
	"syscall"
)

type Runner interface {
	Run() error
}

type Server interface {
	Init() error
	Start() error
	Wait() error
	Close() error
}

type Option interface{}

type defaultRunner struct {
	server Server
}

func NewDefaultRunner(server Server) defaultRunner {
	return defaultRunner{server: server}
}
func (r *defaultRunner) Run() error {
	if err := r.server.Init(); err != nil {
		return err
	}

	if err := r.server.Start(); err != nil {
		return err
	}

	return r.server.Wait()
}

type defaultServer struct {
	Opt Option
}

func (s *defaultServer) Init() error {
	return nil
}

func (s *defaultServer) Wait() error {
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for sign := range signalCh {
		switch sign {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
			s.Close()
			os.Exit(0)
		}
	}

	return nil
}

func (s *defaultServer) Close() error {
	return nil
}

func handlePolicyCheck(policy Policy, events []*event.Event) bool {
	riskLevelFilter := make(map[string]struct{})
	for _, level := range policy.RiskLevelFilter {
		riskLevelFilter[level] = struct{}{}
	}

	for _, event := range events {
		if _, ok := riskLevelFilter[toLevelStr(event.Level)]; !ok {
			continue
		}

		return false
	}

	return true
}
