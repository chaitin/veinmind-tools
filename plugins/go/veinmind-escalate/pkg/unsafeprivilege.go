package pkg

import (
	"bufio"
	"fmt"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
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
var UnSafeCapList = []string{"DAC_OVERRIDE", "DAC_READ_SEARCH", "SYS_MODULE", "SYS_PTRACE", "PRIVILEGED", "SYS_ADMIN"}

func UnsafePrivCheck(fs api.FileSystem) error {
	taskMap := make(map[checkMode][]string, 0)
	taskMap[WRITE] = []string{"/etc/passwd", "/etc/crontab"}
	taskMap[READ] = []string{"/etc/shadow"}

	content, err := fs.Open("/etc/crontab")
	defer content.Close()
	if err == nil {
		scanner := bufio.NewScanner(content)
		res := make([]string, 0)
		for scanner.Scan() {
			if !strings.HasPrefix(scanner.Text(), "#") {
				res = append(res, strings.Split(scanner.Text(), " ")...)
			}
		}
		regexPattern := "^\\/(\\w+\\/?)+.*"
		for _, value := range res {
			if ok, err := regexp.Match(regexPattern, []byte(value)); err == nil {
				if ok == true {
					taskMap[WRITE] = append(taskMap[WRITE], CRONFLAG+value)
				}
			}
		}
	}
	for _, task := range taskMap[WRITE] {
		reason := WRITEREASON
		if strings.HasPrefix(task, CRONFLAG) {
			reason = CRONWRITEREASON
			task = strings.TrimPrefix(task, CRONFLAG)
		}
		if priv, ok, err := PrivCheck(fs, task, WRITE); err == nil {
			if ok == true {
				AddResult(task, reason, "UnSafePriv "+priv)
			}
		} else {
			log.Error(err)
			return err
		}
	}

	for _, task := range taskMap[READ] {
		if priv, ok, err := PrivCheck(fs, task, READ); err == nil {
			if ok == true {
				AddResult(task, READREASON, "UnSafePriv "+priv)
			}
		} else {
			log.Error(err)
			return err
		}
	}
	return nil
}

func UnsafeSuidCheck(fs api.FileSystem) error {
	var binaryName = []string{"bash", "nmap", "vim", "find", "more", "less", "nano", "cp", "awk"}
	var filepath = []string{"/bin/", "/usr/bin/"}
	for i := 0; i < len(filepath); i++ {
		for j := 0; j < len(binaryName); j++ {
			files := filepath[i] + binaryName[j]
			content, err := fs.Stat(files)
			if err == nil {
				if IsBelongToRoot(content) && IsContainSUID(content) {
					AddResult(files, SUIDREASON, "UnSafePriv "+content.Mode().String())
				}
			} else {
				log.Error(err)
				continue
			}
		}
	}
	return nil
}

func UnsafeCapCheck(fs api.FileSystem) error {
	res, err := ReadProc(fs, "/proc/1/status")
	if err != nil {
		log.Error(err)
		return err
	}
	Cap, err := ParseCapEff(res)
	if err != nil {
		log.Error(err)
		return err
	}
	UnSafeCaps := intersect(Cap, UnSafeCapList)
	for _, UnSafeCap := range UnSafeCaps {
		AddResult("/proc/1/status", CAPREASON, "UnSafeCapability "+UnSafeCap)
	}
	return nil
}

func CheckEmptyPasswdRoot(fs api.FileSystem) error {
	privilegedUser := make([]string, 0)
	filePasswd, err := fs.Open("/etc/passwd")
	if err != nil {
		log.Error(err)
		return err
	}
	defer filePasswd.Close()
	if err == nil {
		scanner := bufio.NewScanner(filePasswd)
		for scanner.Scan() {
			attr := strings.Split(scanner.Text(), ":")
			if attr[1] == "" && attr[2] == "0" {
				privilegedUser = append(privilegedUser, attr[0])
			}
		}
	}

	fileShadow, err := fs.Open("/etc/shadow")
	if err != nil {
		log.Error(err)
		return err
	}
	defer fileShadow.Close()
	if err == nil {
		scanner := bufio.NewScanner(fileShadow)
		for scanner.Scan() {
			attr := strings.Split(scanner.Text(), ":")
			if attr[1] == "0" {
				for _, root := range privilegedUser {
					if root == attr[0] {
						AddResult("/etc/shadow", EMPTYPASSWDREASON, "UnsafeUser "+attr[0])
					}
				}
			}
		}
	}
	return nil
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

func PrivCheck(fs api.FileSystem, path string, checkMode checkMode) (string, bool, error) {
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
	if privPasswdAllUsers >= int(checkMode) {
		return content.Mode().String(), true, nil
	}
	return "", false, nil
}

func IsBelongToRoot(content os.FileInfo) bool {
	uid := content.Sys().(*syscall.Stat_t).Uid
	if uid == uint32(0) {
		return true
	}
	return false
}

func IsContainSUID(content os.FileInfo) bool {
	res := fmt.Sprintf("%o", uint32(content.Mode()))
	if strings.HasPrefix(res, "40000") {
		return true
	}
	return false
}

func ParseCapEff(capHex string) ([]string, error) {
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

func ReadProc(fs api.FileSystem, path string) (string, error) {
	//读取 /proc/1/status中的Cap相关数据判断权限
	return "", nil
}
