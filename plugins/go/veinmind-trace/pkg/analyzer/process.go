package analyzer

import (
	api "github.com/chaitin/libveinmind/go"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-trace/pkg/security"
)

// ProcAnalyzer 检测容器内异常的进程
//  1. 隐藏进程(mount -o bind方式) -
//  2. 反弹shell的进程 -
//  3. 带有挖矿、黑客工具、可疑进程名的进程
//  4. 包含 Ptrace 的进程
type ProcAnalyzer struct {
	Object api.Container
	//Event    []

	processes map[int32]api.Process
}

func NewProcAnalyzer(container api.Container) {

}

func (pa *ProcAnalyzer) Scan() {
	pa.scanHideProcess()
	for pid, p := range pa.processes {
		pa.scanReverseShell(p, pid)
		pa.scanEvalProcess(p)
		pa.scanPTraceProcess(p)
	}
}

func (pa *ProcAnalyzer) scanHideProcess() {
	if security.IsHideProcess(pa.Object) {
		// todo
	}
}

func (pa *ProcAnalyzer) scanReverseShell(p api.Process, pid int32) {
	cmdLine, err := p.Cmdline()
	if err != nil {
		return
	}
	if security.IsReverseShell(pa.Object, pid, cmdLine) {
		// todo
	}
}

func (pa *ProcAnalyzer) scanEvalProcess(p api.Process) {
	cmdLine, err := p.Cmdline()
	if err != nil {
		return
	}
	if security.IsEval(cmdLine) {
		// todo
	}
}

func (pa *ProcAnalyzer) scanPTraceProcess(p api.Process) {
	status, err := p.Status()
	if err != nil {
		return
	}
	if security.HasPtraceProcess(status) {
		// todo
	}
}
