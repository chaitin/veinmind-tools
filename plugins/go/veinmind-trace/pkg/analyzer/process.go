package analyzer

import (
	"io"
	"path/filepath"
	"strconv"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report/event"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-trace/pkg/security"
)

func init() {
	Group = append(Group, &ProcAnalyzer{})
}

// ProcAnalyzer 检测容器内异常的进程
//  1. 隐藏进程(mount -o bind方式) -
//  2. 反弹shell的进程 -
//  3. 带有挖矿、黑客工具、可疑进程名的进程
//  4. 包含 Ptrace 的进程
type ProcAnalyzer struct {
	event     []*event.TraceEvent
	container api.Container
}

func (pa *ProcAnalyzer) Scan(container api.Container) {
	pa.event = make([]*event.TraceEvent, 0)
	pa.container = container
	pa.scanHideProcess()
	pids, err := container.Pids()
	if err != nil {
		return
	}

	for _, pid := range pids {
		ps, err := container.NewProcess(pid)
		if err != nil {
			continue
		}
		pa.scanReverseShell(ps, pid)
		pa.scanEvalProcess(ps, pid)
		pa.scanPTraceProcess(container, ps, pid)
	}
}

func (pa *ProcAnalyzer) scanHideProcess() {
	if ok, content := security.IsHideProcess(pa.container); ok {
		pa.event = append(pa.event, &event.TraceEvent{
			Name:        "Hiding Process",
			From:        "Process",
			Path:        "/proc/mounts",
			Description: "some hiding process is in /proc/mounts",
			Detail:      content,
			Level:       event.High,
		})
	}
}

func (pa *ProcAnalyzer) scanReverseShell(p api.Process, pid int32) {
	cmdLine, err := p.Cmdline()
	if err != nil {
		return
	}
	if security.IsReverseShell(pa.container, pid, cmdLine) {
		pa.event = append(pa.event, &event.TraceEvent{
			Name:        "Reverse Shell Process",
			From:        "Process",
			Path:        "/proc/" + strconv.Itoa(int(pid)),
			Description: "an reverse shell process detect",
			Detail:      cmdLine,
			Level:       event.Critical,
		})
	}
}

func (pa *ProcAnalyzer) scanEvalProcess(p api.Process, pid int32) {
	cmdLine, err := p.Cmdline()
	if err != nil {
		return
	}
	if security.IsEval(cmdLine) {
		pa.event = append(pa.event, &event.TraceEvent{
			Name:        "Eval Process",
			From:        "Process",
			Path:        "/proc/" + strconv.Itoa(int(pid)),
			Description: "an eval shell process detect",
			Detail:      cmdLine,
			Level:       event.Critical,
		})
	}
}

func (pa *ProcAnalyzer) scanPTraceProcess(container api.Container, p api.Process, pid int32) {
	cmdLine, err := p.Cmdline()
	if err != nil {
		return
	}

	file, err := container.Open(filepath.Join("/proc", strconv.Itoa(int(pid)), "status"))
	if err != nil {
		return
	}

	status, err := io.ReadAll(file)
	if err != nil {
		return
	}
	if security.HasPtraceProcess(string(status)) {
		pa.event = append(pa.event, &event.TraceEvent{
			Name:        "Ptrace Process",
			From:        "Process",
			Path:        "/proc/" + strconv.Itoa(int(pid)),
			Description: "an process with Ptrace detect",
			Detail:      cmdLine,
			Level:       event.High,
		})
	}
}

func (pa *ProcAnalyzer) Result() []*event.TraceEvent {
	return pa.event
}
