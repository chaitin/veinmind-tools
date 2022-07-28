package plugind

import (
	"context"
	_ "embed"
	"errors"

	"github.com/chaitin/libveinmind/go/plugin/log"
)

// RunService start the Plugin's all Services runner
func RunService(ctx context.Context, runners []*Runner) error {
	err := startAllService(runners)
	if err != nil {
		return err
	}
	daemon(ctx)
	return nil
}

//daemon will create a coroutine user to monitor
//whether the already running service exits.
//If it exits, it will be pulled up again.
func daemon(ctx context.Context) {
	go func() {
		defer kill()
		for {
			select {
			case <-ctx.Done():
				return
			case v := <-Signal:
				err := restart(v)
				if err != nil {
					log.Error(err)
					return
				}
			}
		}
	}()
}

// startAllService start all the runner
func startAllService(runners []*Runner) error {
	for _, runner := range runners {
		RunnerMap.Store(runner.Uuid, runner)
		err := runner.start()
		if err != nil {
			kill()
			return err
		}
	}
	return nil
}

// restart the process according to the Uuid of the process
func restart(uuid string) error {
	value, ok := RunnerMap.Load(uuid)
	if ok {
		return value.(*Runner).start()
	}
	return errors.New("can not find service")
}

// kill all surviving processes
func kill() {
	RunnerMap.Range(func(key, value any) bool {
		if value.(*Runner).Cmd != nil {
			err := value.(*Runner).Cmd.Process.Kill()
			if err != nil {
				log.Error(err)
				return false
			}
		}
		return true
	})
}
