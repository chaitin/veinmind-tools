package plugind

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type ServiceRunner struct {
	ExecScript string
	ExecArgs   []string
	Stderr     *os.File
	Stdout     *os.File
	Port       string
	stop       func()
	Process    *os.Process
	Mut        sync.Mutex
	TimeOut    int
}

func NewService(s ServiceConf) (*ServiceRunner, error) {
	StderrFile, err := os.OpenFile(s.StderrLog, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	StdoutFile, err := os.OpenFile(s.StdoutLog, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &ServiceRunner{
		ExecScript: s.ExecScript,
		ExecArgs:   strings.Split(s.ExecArgs, " "),
		Stdout:     StdoutFile,
		Stderr:     StderrFile,
		Port:       s.Port,
		TimeOut:    10,
	}, nil
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
	sig := make(chan os.Signal, 1)
	go func() {
		signal.Notify(sig)
		signal.Notify(sig, syscall.SIGCHLD)
		cmd := exec.Command(s.ExecScript, s.ExecArgs...) //nolint:gosec
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		cmd.Stderr = s.Stderr
		cmd.Stdout = s.Stdout
		err := cmd.Start()
		if err != nil {
			errs <- fmt.Errorf("panic with error %v", err)
		}
		s.Mut.Lock()
		s.TimeOut = 10
		s.Process = cmd.Process
		s.Mut.Unlock()
		err = cmd.Wait()
		if err != nil {
			errs <- fmt.Errorf("panic with error %v", err)
		}
		sigInfo := <-sig
		errs <- errors.New(sigInfo.String())
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

func (s *ServiceRunner) Ready() error {
	if s.Port != "" {
		for {
			if s.TimeOut <= 0 {
				return errors.New("time out")
			}
			if s.CheckPort() {
				return nil
			} else {
				s.Mut.Lock()
				s.TimeOut = s.TimeOut - 1
				s.Mut.Unlock()
			}
			time.Sleep(time.Duration(s.TimeOut) * time.Second)
		}
	}
	return nil
}

func (s *ServiceRunner) CheckPort() bool {
	if s.Port != "" {
		timeout, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", s.Port), 3*time.Second)
		if err != nil {
			return false
		}
		if timeout != nil {
			return true
		}
		return false
	}
	return true
}
