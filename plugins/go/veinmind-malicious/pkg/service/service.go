package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type ServiceRunner struct {
	ExecScript string
	ExecArgs   []string
	stop       func()
	Process    *os.Process
	mtx        sync.Mutex
}

func NewService(clamavExec, clamavConf string) *ServiceRunner {
	return &ServiceRunner{ExecScript: clamavExec, ExecArgs: strings.Split(clamavConf, " ")}
}

func (s *ServiceRunner) Run(ctx context.Context) {
	ctx, s.stop = context.WithCancel(ctx)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-s.loop(ctx):
			}
		}
	}()
}

func (s *ServiceRunner) loop(ctx context.Context) <-chan error {
	errs := make(chan error)
	go func() {
		clam := exec.Command(s.ExecScript, s.ExecArgs...) //nolint:gosec
		err := clam.Start()
		if err != nil {
			errs <- fmt.Errorf("panic with error %v", err)
		}
		s.mtx.Lock()
		s.Process = clam.Process
		s.mtx.Unlock()
		err = clam.Wait()
		if err != nil {
			errs <- fmt.Errorf("panic with error %v", err)
		}
		errs <- errors.New("service exit")
	}()
	return errs
}

func (s *ServiceRunner) Stop() error {
	if s.stop != nil {
		s.stop()
		if s.Process != nil {
			err := s.Process.Kill()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
