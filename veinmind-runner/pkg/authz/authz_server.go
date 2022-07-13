package authz

import (
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type ServerOption func(option *serverOption) error

type serverOption struct {
	authLog   *os.File
	pluginLog *os.File
	policies  *policyMap
	listener  net.Listener
}

func newServerOption(options ...ServerOption) (*serverOption, error) {
	result := &serverOption{
		policies: newPolicyMap(),
	}
	for _, option := range options {
		if err := option(result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func WithPolicy(policies ...Policy) ServerOption {
	return func(option *serverOption) error {
		if option.policies == nil {
			option.policies = newPolicyMap()
		}
		for _, policy := range policies {
			option.policies.Store(policy)
		}

		return nil
	}
}

func WithAuthLog(path string) ServerOption {
	return func(option *serverOption) error {
		_, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			_, _ = os.Create(path)
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
			_, _ = os.Create(path)
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

func WithServerOptions(options []ServerOption) ServerOption {
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

func (my *defaultServer) Run() error {
	defer func() {
		var result interface{} = my.server
		option, _ := result.(serverOption)
		_ = option.authLog.Close()
		_ = option.pluginLog.Close()
	}()

	if err := my.init(); err != nil {
		return err
	}

	if err := my.start(); err != nil {
		return err
	}

	return my.wait()
}

func (my *defaultServer) init() error {
	var result *serverOption
	switch s := my.server.(type) {
	case *dockerPluginServer:
		result = (*serverOption)(s)
	default:
		return errors.New("not support the server")
	}
	if err := WithServerOptions(my.options)(result); err != nil {
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
	if err := WithServerOptions(defaultOptions)(result); err != nil {
		return err
	}

	return my.server.init()
}

func (my *defaultServer) wait() error {
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

	for s := range signalCh {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
			os.Exit(0)
		case syscall.SIGUSR1:
		case syscall.SIGUSR2:
		default:
		}
	}

	return my.server.wait()
}
