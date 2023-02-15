package pkg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/chaitin/libveinmind/go/pkg/vfs"
	"github.com/chaitin/veinmind-common-go/service/report/event"
	"os"
	"strconv"
	"strings"
	"syscall"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/plugin/log"
)

type checkMode int

var hostConfig = make(map[string]interface{}, 0)
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
	err := getContainerHostConfig(fs)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var res = make([]*event.EscapeDetail, 0)
	if hostConfig["Privileged"] == true {
		res = append(res, &event.EscapeDetail{
			Target: "LINUX CAPABILITY",
			Reason: CAPREASON,
			Detail: "UnSafeCapability PRIVILEGED",
		})
	} else {
		Caps := hostConfig["CapAdd"]
		if Caps != nil {
			UnSafeCaps := intersect(UnSafeCapList, Caps.([]interface{}))
			for _, value := range UnSafeCaps {
				if value == "CAP_SYS_PTRACE" {
					if hostConfig["PidMode"] == "host" {
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

func getContainerHostConfig(container api.FileSystem) error {
	fileContent, err := vfs.Open("/var/lib/docker/containers/" + strings.TrimPrefix(container.(api.Container).ID(), "sha256:") + "/hostconfig.json")
	if err != nil {
		return err
	}
	defer fileContent.Close()
	fileInfo, err := vfs.Stat("/var/lib/docker/containers/" + strings.TrimPrefix(container.(api.Container).ID(), "sha256:") + "/hostconfig.json")
	if err != nil {
		return err
	}
	content := make([]byte, fileInfo.Size())
	fileContent.Read(content)
	err = json.Unmarshal(content, &hostConfig)
	if err != nil {
		return err
	}
	return nil
}

func intersect(a []string, b []interface{}) []string {
	iter := make([]string, 0)
	mp := make(map[string]bool, 0)
	for _, value := range a {
		if _, ok := mp[value]; !ok {
			mp[value] = true
		}
	}
	for _, value := range b {
		str, err := value.(string)
		if err == false {
			log.Error(value, "is not string")
			return nil
		}
		if _, ok := mp[str]; ok {
			iter = append(iter, str)
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

func init() {
	ContainerCheckList = append(ContainerCheckList, ContainerUnsafeCapCheck, UnsafePrivCheck, UnsafeSuidCheck, CheckEmptyPasswdRoot)
	ImageCheckList = append(ImageCheckList, UnsafePrivCheck, UnsafeSuidCheck, CheckEmptyPasswdRoot)
}
