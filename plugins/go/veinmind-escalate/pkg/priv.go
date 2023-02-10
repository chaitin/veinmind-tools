package pkg

import (
	"bufio"
	"fmt"
	"github.com/chaitin/veinmind-common-go/service/report/event"
	"os"
	"strconv"
	"strings"
	"syscall"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/plugin/log"
)

type checkMode int

var CAPStringsList = []string{
	"CAP_CHOWN",
	"CAP_DAC_OVERRIDE",
	"CAP_DAC_READ_SEARCH",
	"CAP_FOWNER",
	"CAP_FSETID",
	"CAP_KILL",
	"CAP_SETGID",
	"CAP_SETUID",
	"CAP_SETPCAP",
	"CAP_LINUX_IMMUTABLE",
	"CAP_NET_BIND_SERVICE",
	"CAP_NET_BROADCAST",
	"CAP_NET_ADMIN",
	"CAP_NET_RAW",
	"CAP_IPC_LOCK",
	"CAP_IPC_OWNER",
	"CAP_SYS_MODULE",
	"CAP_SYS_RAWIO",
	"CAP_SYS_CHROOT",
	"CAP_SYS_PTRACE",
	"CAP_SYS_PACCT",
	"CAP_SYS_ADMIN",
	"CAP_SYS_BOOT",
	"CAP_SYS_NICE",
	"CAP_SYS_RESOURCE",
	"CAP_SYS_TIME",
	"CAP_SYS_TTY_CONFIG",
	"CAP_MKNOD",
	"CAP_LEASE",
	"CAP_AUDIT_WRITE",
	"CAP_AUDIT_CONTROL",
	"CAP_SETFCAP",
	"CAP_MAC_OVERRIDE",
	"CAP_MAC_ADMIN",
	"CAP_SYSLOG",
	"CAP_WAKE_ALARM",
	"CAP_BLOCK_SUSPEND",
	"CAP_AUDIT_READ",
	"CAP_PERFMON",
	"CAP_BPF",
	"CAP_CHECKPOINT_RESTORE",
}
var UnSafeCapList = []string{"CAP_DAC_READ_SEARCH", "CAP_SYS_MODULE", "CAP_SYS_PTRACE", "PRIVILEGED", "CAP_SYS_ADMIN", "CAP_SYS_CHROOT", "CAP_BPF", "CAP_DAC_OVERRIDE"}

func UnsafePrivCheck(fs api.FileSystem) ([]*event.EscapeDetail, error) {
	var res = make([]*event.EscapeDetail, 0)
	taskMap := make(map[checkMode][]string, 0)
	taskMap[WRITE] = []string{"/etc/passwd", "/etc/crontab"}
	taskMap[READ] = []string{"/etc/shadow"}
	for _, task := range taskMap[WRITE] {
		reason := WRITEREASON
		if priv, ok, err := privCheck(fs, task, WRITE); err == nil {
			if ok == true {
				res = append(res, &event.EscapeDetail{
					Target: task,
					Reason: reason,
					Detail: "UnSafePriv " + priv,
				})
			}
		}
	}
	for _, task := range taskMap[READ] {
		if priv, ok, err := privCheck(fs, task, READ); err == nil {
			if ok == true {
				res = append(res, &event.EscapeDetail{
					Target: task,
					Reason: READREASON,
					Detail: "UnSafePriv " + priv,
				})
			}
		}
	}
	return res, nil
}

func UnsafeSuidCheck(fs api.FileSystem) ([]*event.EscapeDetail, error) {
	var res = make([]*event.EscapeDetail, 0)
	var binaryName = []string{"bash", "nmap", "vim", "find", "more", "less", "nano", "cp", "awk"}
	var filepath = []string{"/bin/", "/usr/bin/"}
	for i := 0; i < len(filepath); i++ {
		for j := 0; j < len(binaryName); j++ {
			files := filepath[i] + binaryName[j]
			content, err := fs.Stat(files)
			if err == nil {
				if isBelongToRoot(content) && isContainSUID(content) {
					res = append(res, &event.EscapeDetail{
						Target: files,
						Reason: SUIDREASON,
						Detail: "UnSafePriv " + content.Mode().String(),
					})
				}
			} else {
				continue
			}
		}
	}
	return res, nil
}

func ContainerUnsafeCapCheck(fs api.FileSystem) ([]*event.EscapeDetail, error) {
	var res = make([]*event.EscapeDetail, 0)
	file, err := fs.Open("/proc/1/status")
	if err != nil {
		log.Error(err)
		return res, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "CapEff:") {
			if strings.HasSuffix(scanner.Text(), "fffffffff") {
				res = append(res, &event.EscapeDetail{
					Target: "/proc/1/status",
					Reason: CAPREASON,
					Detail: "UnSafeCapability PRIVILEGED",
				})
			} else {
				Cap, err := parseCapEff(scanner.Text())
				if err != nil {
					log.Error(err)
					return res, err
				}
				UnSafeCaps := intersect(Cap, UnSafeCapList)
				for _, UnSafeCap := range UnSafeCaps {
					res = append(res, &event.EscapeDetail{
						Target: "/proc/1/status",
						Reason: CAPREASON,
						Detail: "UnSafeCapability " + UnSafeCap,
					})
				}
			}
		}
	}
	return res, nil
}

func CheckEmptyPasswdRoot(fs api.FileSystem) ([]*event.EscapeDetail, error) {
	var res = make([]*event.EscapeDetail, 0)
	privilegedUser := make([]string, 0)
	filePasswd, err := fs.Open("/etc/passwd")
	if err != nil {
		return res, err
	}
	defer filePasswd.Close()
	scanner := bufio.NewScanner(filePasswd)
	for scanner.Scan() {
		attr := strings.Split(scanner.Text(), ":")
		if len(attr) >= 3 {
			if attr[2] == "0" {
				privilegedUser = append(privilegedUser, attr[0])
			}
		}
	}

	fileShadow, err := fs.Open("/etc/shadow")
	if err != nil {
		return res, err
	}
	defer fileShadow.Close()
	scanner = bufio.NewScanner(fileShadow)
	for scanner.Scan() {
		attr := strings.Split(scanner.Text(), ":")
		if attr[1] == "0" {
			for _, root := range privilegedUser {
				if root == attr[0] {
					res = append(res, &event.EscapeDetail{
						Target: "/etc/shadow",
						Reason: EMPTYPASSWDREASON,
						Detail: "UnsafeUser " + attr[0],
					})
				}
			}
		}
	}

	return res, nil
}

func intersect(a []string, b []string) []string {
	iter := make([]string, 0)
	mp := make(map[string]bool, 0)
	for _, value := range a {
		if _, ok := mp[value]; !ok {
			mp[value] = true
		}
	}
	for _, value := range b {
		if _, ok := mp[value]; ok {
			iter = append(iter, value)
		}
	}
	return iter
}

func privCheck(fs api.FileSystem, path string, checkMode checkMode) (string, bool, error) {
	content, err := fs.Stat(path)
	if err != nil {
		return "", false, err
	}
	mode := fmt.Sprintf("%o", uint32(content.Mode()))
	privPasswdAllUsers, err := strconv.Atoi(string(mode[len(mode)-1]))
	if err != nil {
		log.Error(err)
		return "", false, err
	}
	if checkMode == WRITE {
		if privPasswdAllUsers >= int(checkMode) && privPasswdAllUsers != 4 {
			return content.Mode().String(), true, nil
		}
	} else {
		if privPasswdAllUsers >= int(checkMode) {
			return content.Mode().String(), true, nil
		}
	}
	return "", false, nil
}

func isBelongToRoot(content os.FileInfo) bool {
	uid := content.Sys().(*syscall.Stat_t).Uid
	if uid == uint32(0) {
		return true
	}
	return false
}

func isContainSUID(content os.FileInfo) bool {
	res := fmt.Sprintf("%o", uint32(content.Mode()))
	if strings.HasPrefix(res, "40000") {
		return true
	}
	return false
}

func parseCapEff(capHex string) ([]string, error) {
	capHex = strings.TrimSpace(strings.TrimPrefix(capHex, "CapEff:"))
	var capTextList []string
	numb, err := strconv.ParseUint(capHex, 16, 64)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(CAPStringsList); i++ {
		flag := numb & 0x1
		if flag == uint64(1) {
			capTextList = append(capTextList, CAPStringsList[i])
		}
		numb = numb >> 1
	}

	return capTextList, nil
}

func init() {
	ContainerCheckList = append(ContainerCheckList, ContainerUnsafeCapCheck, UnsafePrivCheck, UnsafeSuidCheck, CheckEmptyPasswdRoot)
	ImageCheckList = append(ImageCheckList, UnsafePrivCheck, UnsafeSuidCheck, CheckEmptyPasswdRoot)
}
