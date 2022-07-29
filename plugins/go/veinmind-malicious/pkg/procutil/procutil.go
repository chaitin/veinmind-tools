package procutil

import (
	"errors"
	ps "github.com/shirou/gopsutil/process"
)

func GetProcessByPID(PID int) (*ps.Process, error) {
	processes, err := ps.Processes()
	if err != nil {
		return nil, err
	}
	for _, process := range processes {
		if err != nil {
			return nil, err
		}
		if int(process.Pid) == PID {
			return process, nil
		}
	}
	return nil, errors.New("process of PID not Found")
}
