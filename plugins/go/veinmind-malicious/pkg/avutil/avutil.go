package avutil

import (
	"context"
	"errors"
	"net"
	"os/exec"
	"sync"
	"time"
)

type ServiceOption func(*ClamAVManager)

func WithPort(port string) ServiceOption {
	return func(manger *ClamAVManager) {
		manger.clamAVPort = port
	}
}

func WithExec(exec string) ServiceOption {
	return func(manger *ClamAVManager) {
		manger.clamAVExec = exec
	}
}

func WithConf(config string) ServiceOption {
	return func(manger *ClamAVManager) {
		manger.clamAVConf = config
	}
}

func WithHost(host string) ServiceOption {
	return func(manger *ClamAVManager) {
		manger.clamAVHost = host
	}
}

type ClamAVManager struct {
	ctx        context.Context
	sig        chan struct{}
	wg         *sync.WaitGroup
	clamAVPort string
	clamAVHost string
	clamAVExec string
	clamAVConf string
}

func New(ctx context.Context, opts ...ServiceOption) *ClamAVManager {
	cam := &ClamAVManager{
		ctx: ctx,
		sig: make(chan struct{}),
		wg:  &sync.WaitGroup{},
	}
	for _, opt := range opts {
		opt(cam)
	}
	return cam
}

func (c *ClamAVManager) Run() error {
	clamAVRunner := exec.CommandContext(c.ctx, c.clamAVExec, "-c", c.clamAVConf, "-F")
	err := clamAVRunner.Run()
	if err != nil {
		return err
	}
	c.sig <- struct{}{}
	return nil
}

func (c *ClamAVManager) Ready() error {
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

func (c *ClamAVManager) checkPortIsUsed() error {
	_, err := net.DialTimeout("tcp", net.JoinHostPort(c.clamAVHost, c.clamAVPort), 3*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClamAVManager) Daemon() error {
	c.wg.Add(1)
	defer c.wg.Done()
	for {
		select {
		case <-c.ctx.Done():
			return nil
		case <-c.sig:
			err := c.Run()
			if err != nil {
				return err
			}
		}
	}
}

func (c *ClamAVManager) Wait() {
	c.wg.Wait()
}
