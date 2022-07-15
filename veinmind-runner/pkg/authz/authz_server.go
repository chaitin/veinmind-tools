package authz

import (
	"errors"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/reporter"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type ServerOption func(option *serverOption) error

type serverOption struct {
	authLog   io.WriteCloser
	pluginLog io.WriteCloser
	policies  sync.Map
	listener  net.Listener
}

func WithPolicy(policies ...Policy) ServerOption {
	return func(option *serverOption) error {
		for _, policy := range policies {
			option.policies.Store(policy.Action, policy)
		}

		return nil
	}
}

func WithAuthLog(path string) ServerOption {
	return func(option *serverOption) error {
		_, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			_, err = os.Create(path)
			if err != nil {
				return err
			}
		}

		fp, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}

		option.authLog = fp
		return nil
	}
}

func WithPluginLog(path string) ServerOption {
	return func(option *serverOption) error {
		_, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			_, err = os.Create(path)
			if err != nil {
				return err
			}
		}

		fp, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}

		option.pluginLog = fp
		return nil
	}
}

func WithListenerUnix(addr string) ServerOption {
	return func(option *serverOption) error {
		listener, err := net.ListenUnix("unix", &net.UnixAddr{Net: "unix", Name: addr})
		if err != nil {
			return err
		}

		option.listener = listener
		return nil
	}
}

func WithServerOptions(options ...ServerOption) ServerOption {
	return func(s *serverOption) error {
		for _, option := range options {
			if err := option(s); err != nil {
				return err
			}
		}

		return nil
	}
}

type Server interface {
	Run() error
}

type server interface {
	init() error
	start() (err error)
	wait() error
	close()
}

type defaultServer struct {
	server
	options []ServerOption
}

func newDefaultServer(s server, opts ...ServerOption) *defaultServer {
	result := &defaultServer{
		server:  s,
		options: make([]ServerOption, 0),
	}
	result.options = append(result.options, opts...)

	return result
}

func (s *defaultServer) Run() error {
	if err := s.init(); err != nil {
		return err
	}

	if err := s.start(); err != nil {
		return err
	}

	return s.wait()
}

func (s *defaultServer) init() error {
	var result *serverOption
	switch srv := s.server.(type) {
	case *dockerPluginServer:
		result = srv.option
	default:
		return errors.New("not support the server")
	}
	if err := WithServerOptions(s.options...)(result); err != nil {
		return err
	}

	var defaultOptions []ServerOption
	if result.authLog == nil {
		defaultOptions = append(defaultOptions, WithAuthLog(defaultAuthLogPath))
	}
	if result.pluginLog == nil {
		defaultOptions = append(defaultOptions, WithPluginLog(defaultPluginPath))
	}
	if result.listener == nil {
		defaultOptions = append(defaultOptions, WithListenerUnix(defaultSockListenAddr))
	}
	if err := WithServerOptions(defaultOptions...)(result); err != nil {
		return err
	}

	return s.server.init()
}

func (s *defaultServer) wait() error {
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for sign := range signalCh {
		switch sign {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
			s.close()
			os.Exit(0)
		}
	}

	return nil
}

func (s *defaultServer) close() {
	s.server.close()
}

func handlePolicyCheck(policy Policy, events []reporter.ReportEvent) bool {
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
