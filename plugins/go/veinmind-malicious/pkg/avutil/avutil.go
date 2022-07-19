package avutil

import (
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/pkg/netutil"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/pkg/procutil"
	"github.com/pkg/errors"
	"os/exec"
)

func StartClamAV(port, clamavExec, clamavConf string) (int, error) {
	clam := exec.Command(clamavExec, "-c", clamavConf) //nolint:gos-ec
	err := clam.Run()
	if err != nil {
		return 0, err
	}
	PID, err := CheckClamAVPortStatus(port)
	if err != nil {
		return 0, err
	}
	if PID != 0 {
		return PID, nil
	}
	return 0, errors.New("failed to activate clamAV")
}

func CloseClamAV(PID int) error {
	ClamavProcess, err := procutil.GetProcessByPID(PID)
	if err != nil {
		return err
	}
	name, err := ClamavProcess.Name()
	if err != nil {
		return err
	}
	if name == "clamd" {
		err := ClamavProcess.Kill()
		if err != nil {
			return err
		}
	} else {
		return errors.New("not clamAV process")
	}
	return nil
}

// ClamAVPreCheck check if the clamAV is running and the host is local
func ClamAVPreCheck(host, port string) error {
	local, err := netutil.IsLocalHost(host)
	if err != nil {
		return err
	}
	if !local {
		return errors.New("not local host")
	}
	PID, err := CheckClamAVPortStatus(port)
	if err != nil {
		return err
	}
	if PID != 0 {
		return errors.New("port occupation")
	}
	return nil
}

// checkClamAVPortStatus check if the port of ClamAV is opening
// err is nil and pid is 0: can not find the pid
func CheckClamAVPortStatus(port string) (int, error) {
	PID, err := netutil.GetPidByPort(port)
	if err != nil {
		return 0, err
	}
	if PID != 0 {
		return PID, nil
	}
	return 0, nil
}
