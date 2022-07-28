package plugind

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"os"
	"syscall"
	"time"
)

func newRunner(s ServiceConf) (*Runner, error) {
	StderrFile, err := os.OpenFile(s.StderrLog, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	StdoutFile, err := os.OpenFile(s.StdoutLog, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &Runner{
		Name:    s.Name,
		Uuid:    uuid.New().String(),
		Command: s.Command,
		Stdout:  StdoutFile,
		Stderr:  StderrFile,
		Port:    s.Port,
		TimeOut: 10 * time.Second,
	}, nil
}

//run the Service
//time out error:The target port has not been listened to for a certain period of time
func (s *Runner) start() error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)

	err := s.run()
	if err != nil {
		return err
	}
	if s.Port == "" {
		return nil
	}

	g.Go(func() error {
		for {
			time.Sleep(time.Second)
			select {
			case <-ctx.Done():
				return errors.New("time out")
			default:
				if s.CheckPort() {
					return nil
				}
			}
		}
	})

	return g.Wait()
}

func (s *Runner) CheckPort() bool {
	timeout, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", s.Port), 3*time.Second)
	if err != nil {
		return false
	}
	if timeout != nil {
		defer timeout.Close()
		return true
	}
	return false
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
			log.Println(s.Name, ":", err)
		}
		Signal <- s.Uuid
	}()
	s.Cmd = cmd

	return nil
}
