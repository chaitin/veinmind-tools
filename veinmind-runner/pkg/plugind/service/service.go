package service

import (
	"context"
	"errors"
	"github.com/chaitin/libveinmind/go/plugin/log"
	ps "github.com/shirou/gopsutil/process"
	"golang.org/x/sync/errgroup"
	"syscall"
	"time"
)

// Start run the Service
//time out error:The target port has not been ready to for a certain period of time
func (s *Runner) Start() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.TimeOut)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)

	err := s.run()
	if err != nil {
		return err
	}

	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return errors.New("time out")
			case <-time.Tick(time.Second):
				stat, err := s.Status()
				if err != nil {
					return err
				}
				if stat == "S" {
					return nil
				}
			}
		}
	})

	return g.Wait()
}

func (s *Runner) run() error {
	cmd, err := createCommand(s.Command)
	if err != nil {
		return err
	}
	cmd.Stderr = s.Stderr
	cmd.Stdout = s.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	err = cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Error(s.Command, ":", err)
		}
		s.Signal <- s.Uuid
	}()
	s.Cmd = cmd

	return nil
}

func (s *Runner) Status() (string, error) {
	if s.Cmd == nil {
		return "", errors.New("process not running")
	}
	if s.Cmd.ProcessState != nil {
		return "", errors.New("process exit")
	}
	process, err := ps.NewProcess(int32(s.Cmd.Process.Pid))
	if err != nil {
		return "", err
	}
	// Status returns the process status.
	// R: Running S: Sleep T: Stop I: Idle
	// Z: Zombie W: Wait L: Lock
	stat, err := process.Status()
	if err != nil {
		return "", err
	}
	return stat, nil
}

func (s *Runner) Stop() error {
	if s.Cmd != nil {
		err := s.Cmd.Process.Kill()
		if err != nil {
			return err
		}
	}
	return nil
}
