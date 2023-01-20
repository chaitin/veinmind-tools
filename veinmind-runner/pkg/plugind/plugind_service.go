package plugind

import (
	"context"
	"errors"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/log"
	"net"
	"os"
	"os/exec"
	"sync"
	"time"
)

type serviceOption func(*service)

func withStdout(stdout string) serviceOption {
	return func(s *service) {
		s.stdout = stdout
	}
}

func withStderr(stderr string) serviceOption {
	return func(s *service) {
		s.stderr = stderr
	}
}

func withCheckChains(checks ...serviceCheckFunc) serviceOption {
	return func(s *service) {
		s.checkChains = append(s.checkChains, checks...)
	}
}

func withTimeout(timeout int) serviceOption {
	return func(s *service) {
		s.timeout = time.Duration(timeout) * time.Second
	}
}

func withWaitGroup(wg *sync.WaitGroup) serviceOption {
	return func(s *service) {
		s.wg = wg
	}
}

type service struct {
	ctx         context.Context
	sig         chan struct{}
	cmd         string
	stdout      string
	stderr      string
	wg          *sync.WaitGroup
	timeout     time.Duration
	checkChains []serviceCheckFunc
}

func (s *service) start() error {
	command := exec.CommandContext(s.ctx, "/bin/bash", "-c", s.cmd)

	stdout, err := os.OpenFile(s.stdout, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	stderr, err := os.OpenFile(s.stderr, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	command.Stdout = stdout
	command.Stderr = stderr

	go func() {
		err := command.Run()
		if err != nil {
			log.GetModule(log.PlugindModuleKey).Error(err)
		}
		s.sig <- struct{}{}
	}()

	return nil
}

func (s *service) ready() error {
	if s.checkChains == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return errors.New("time out")
		case <-time.Tick(time.Second):
			for _, chain := range s.checkChains {
				ok, err := chain(s)
				if err != nil {
					return err
				}
				if ok {
					return nil
				}
			}
		}
	}
}

func (s *service) daemon() {
	s.wg.Add(1)
	defer s.wg.Done()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.sig:
			err := s.start()
			if err != nil {
				log.GetModule(log.PlugindModuleKey).Error(err)
				return
			}
		}
	}
}

func newService(ctx context.Context, cmd string, opts ...serviceOption) *service {
	svc := &service{
		ctx: ctx,
		cmd: cmd,
		sig: make(chan struct{}),
	}

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

type serviceCheckFunc func(*service) (bool, error)

var serviceChecks = map[string]func(string) serviceCheckFunc{
	"file": SVC_CHECK_FILE,
	"port": SVC_CHECK_PORT,
}

var (
	SVC_CHECK_FILE = checkFileIsExisted
	SVC_CHECK_PORT = checkPortIsUsed
)

func checkFileIsExisted(path string) serviceCheckFunc {
	return func(s *service) (bool, error) {
		_, err := os.Stat(path)
		if err == nil {
			return true, nil
		}
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		} else {
			return false, err
		}
	}
}

func checkPortIsUsed(port string) serviceCheckFunc {
	return func(s *service) (bool, error) {
		_, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", port), 3*time.Second)
		if err != nil {
			return false, nil
		}
		return true, nil
	}
}
