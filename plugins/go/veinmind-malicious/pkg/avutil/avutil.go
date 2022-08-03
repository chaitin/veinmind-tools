package avutil

import (
	"context"
	"errors"
	"net"
	"os"
	"os/exec"
	"sync"
	"time"
)

type ServiceOption func(*ClamAVManger)

func WithPort(port string) ServiceOption {
	return func(manger *ClamAVManger) {
		manger.clamAVPort = port
	}
}

func WithExec(exec string) ServiceOption {
	return func(manger *ClamAVManger) {
		manger.clamAVExec = exec
	}
}

func WithConf(config string) ServiceOption {
	return func(manger *ClamAVManger) {
		manger.clamAVConf = config
	}
}

func WithHost(host string) ServiceOption {
	return func(manger *ClamAVManger) {
		manger.clamAVHost = host
	}
}

type ClamAVManger struct {
	ctx        context.Context
	sig        chan struct{}
	wg         *sync.WaitGroup
	clamAVPort string
	clamAVHost string
	clamAVExec string
	clamAVConf string
	proc       *os.Process
}

func New(ctx context.Context, opts ...ServiceOption) *ClamAVManger {
	cam := &ClamAVManger{
		ctx: ctx,
		sig: make(chan struct{}),
		wg:  &sync.WaitGroup{},
	}
	for _, opt := range opts {
		opt(cam)
	}
	return cam
}

func (c *ClamAVManger) Run() error {
	clamAVRunner := exec.Command(c.clamAVExec, "-c", c.clamAVConf, "-F")
	err := clamAVRunner.Start()
	if err != nil {
		return err
	}
	c.proc = clamAVRunner.Process
	err = clamAVRunner.Wait()
	if err != nil {
		return err
	}
	c.sig<- struct{}{}
	return nil
}

func (c *ClamAVManger) Ready() error {
	ctx, cancel := context.WithTimeout(c.ctx, 20*time.Second)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return errors.New("time out")
		case <-time.Tick(time.Second):
			err := c.checkPortIsUsed()
			if err != nil {
				continue
			}
			return nil
		}
	}
}

func (c *ClamAVManger) checkPortIsUsed() error {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(c.clamAVHost, c.clamAVPort), 3*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func (c *ClamAVManger) Daemon() error {
	c.wg.Add(1)
	defer c.wg.Done()
	for {
		select {
		case <-c.ctx.Done():
			if c.proc != nil {
				err := c.proc.Kill()
				if err != nil {
					return err
				}
			}
			return nil
		case <-c.sig:
			err := c.Run()
			if err != nil {
				return err
			}
		}
	}
}

func (c *ClamAVManger) Wait() {
	c.wg.Wait()
}
