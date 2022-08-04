package plugind

import (
	"context"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"net"
	"os"
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

func (s *service) Start(wg *sync.WaitGroup) error {
	err := s.run()
	if err != nil {
		return err
	}

	// check service is working
	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return errors.New("time out")
			case <-time.Tick(time.Second):
				if s.checkChains == nil {
					return nil
				}
				for _, chain := range s.checkChains {
					err := chain(s)
					if err == nil {
						return nil
					}
				}
			}
		}
	})
	if err = g.Wait(); err != nil {
		return err
	}

	// daemon this process
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-s.sig:
				err := s.run()
				if err != nil {
					log.Error(err)
					return
				}
			}
		}
	}()

	return nil
}

func (s *service) run() error {
	command, err := createCommandWithContext(s.ctx, s.cmd)

	if err != nil {
		return err
	}
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
			log.Error(err)
		}
		s.sig <- struct{}{}
	}()

	return nil
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

type serviceCheckFunc func(*service) error

var serviceChecks = map[string]func(string) serviceCheckFunc{
	"file": SVC_CHECK_FILE,
	"port": SVC_CHECK_PORT,
}

var (
	SVC_CHECK_FILE = checkFileIsExisted
	SVC_CHECK_PORT = checkPortIsUsed
)

func checkFileIsExisted(path string) serviceCheckFunc {
	return func(s *service) error {
		_, err := os.Stat(path)
		return err
	}
}

func checkPortIsUsed(port string) serviceCheckFunc {
	return func(s *service) error {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", port), 3*time.Second)
		if err != nil {
			return err
		}
		defer conn.Close()
		return nil
	}
}
