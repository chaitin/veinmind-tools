package pkg

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report/event"
)

type checkMode int

var cap = make([]string, 0)
var UnSafeCapList = []string{"CAP_DAC_READ_SEARCH", "CAP_SYS_MODULE", "CAP_SYS_PTRACE", "CAP_SYS_ADMIN", "CAP_DAC_OVERRIDE"}

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
	container, ok := fs.(api.Container)
	if ok == false {
		log.Error(fs, "is not a container")
		return nil, nil
	}
	err := getCapEff(container)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if isPrivileged(container) {
		res = append(res, &event.EscapeDetail{
			Target: "LINUX CAPABILITY",
			Reason: CAPREASON,
			Detail: "UnSafeCapability PRIVILEGED",
		})
	} else {
		UnSafeCap := intersect(cap, UnSafeCapList)

		for _, value := range UnSafeCap {
			if value == "CAP_SYS_PTRACE" {
				if isPidEqualHost(container) {
					res = append(res, &event.EscapeDetail{
						Target: "LINUX CAPABILITY",
						Reason: CAPREASON,
						Detail: "UnSafeCapability " + value + " and you start the container with parameter : `pid=host`",
					})
				}
			} else {
				res = append(res, &event.EscapeDetail{
					Target: "LINUX CAPABILITY",
					Reason: CAPREASON,
					Detail: "UnSafeCapability " + value,
				})
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

func isPrivileged(container api.Container) bool {
	state, err := container.OCIState()
	if err != nil {
		log.Error(err)
		return false
	}

	if state.Pid == 0 {
		return false
	}

	status, err := os.ReadFile(filepath.Join(func() string {
		fs := os.Getenv("LIBVEINMIND_HOST_ROOTFS")
		if fs == "" {
			return "/"
		}
		return fs
	}(), "proc", strconv.Itoa(state.Pid), "status"))
	if err != nil {
		log.Error(err)
		return false
	}

	pattern := regexp.MustCompile(`(?i)capeff:\s*?([a-z0-9]+)\s`)
	matched := pattern.FindStringSubmatch(string(status))

	if len(matched) != 2 {
		return false
	}

	if strings.HasSuffix(matched[1], "ffffffff") {
		return true
	}

	return false
}

func isPidEqualHost(container api.Container) bool {
	spec, err := container.OCISpec()
	if err != nil {
		return false
	}
	namespaces := spec.Linux.Namespaces
	for _, value := range namespaces {
		if value.Type == "pid" {
			return false
		}
	}
	return true
}

func getCapEff(container api.Container) error {
	spec, err := container.OCISpec()
	if err != nil {
		return err
	}
	cap = append(cap, spec.Process.Capabilities.Effective...)
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
	// 4000 is SUID, 6000 is SUID and SGID
	if strings.HasPrefix(res, "4000") || strings.HasPrefix(res, "6000") {
		return true
	}
	return false
}

func init() {
	ContainerCheckList = append(ContainerCheckList, ContainerUnsafeCapCheck, UnsafePrivCheck, UnsafeSuidCheck, CheckEmptyPasswdRoot)
	ImageCheckList = append(ImageCheckList, UnsafePrivCheck, UnsafeSuidCheck, CheckEmptyPasswdRoot)
}
