package utils

import (
	"bufio"
	"encoding/json"
	_ "encoding/json"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/veinmind-common-go/service/report"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-escalate/models"
	_ "github.com/docker/docker/api/types/mount"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var res = []*models.EscalateResult{}
var escalateLock sync.Mutex
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

//--------------------------容器内提权相关-------------------------------------

// 不安全的suid配置中判断某个可执行文件是否配置了suid

func IsBelongToRoot(content os.FileInfo, s string) bool {

	uid := content.Sys().(*syscall.Stat_t).Uid
	//log.Info(uid.Interface().(uint32))
	if uid == uint32(0) {
		return true
	}
	return false
}

// 不安全的suid配置中判断某个配置了suid权限的可执行文件的属主是否是root

func IsContainSUID(content os.FileInfo, s string) bool {
	res := content.Mode().String()
	if strings.Contains(res, "u") {
		return true
	}
	return false
}

// 不安全的suid配置

func FindSuid(fs api.FileSystem) {
	var binaryName = []string{"bash", "nmap", "vim", "find", "more", "less", "nano", "cp", "awk"}
	var filepath = []string{"/bin/", "/usr/bin/"}
	for i := 0; i < len(filepath); i++ {
		for j := 0; j < len(binaryName); j++ {
			files := filepath[i] + binaryName[j]
			//log.Info(files)
			content, err := fs.Stat(files)
			if err == nil {
				if IsBelongToRoot(content, files) && IsContainSUID(content, files) {
					AddResult(files, "UnSafeSuid")
				}
			} else {
				continue
			}
		}
	}

}

// 空密码高权限用户

func CheckEmptyPasswdRoot(fs api.FileSystem) {
	//更改为shadow

	path := "/etc/passwd"
	file, err := fs.Open(path)
	defer file.Close()
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			attr := strings.Split(scanner.Text(), ":")
			if attr[1] == "" && attr[2] == "0" {
				AddResult(attr[0], "EmptyPasswdRoot")
			}
		}
	}
}

// 检查文件权限

func FileWriteCheck(s string) bool {
	resW := 0
	for i := 0; i < len(s); i++ {
		if s[i] == 'w' {
			resW++
		}
	}
	if resW > 1 {
		return true
	} else {
		return false
	}
}

// 不安全的权限配置，如所有用户可写/etc/passwd，/etc/crontab的定时任务中有某个脚本所有人可写

func UnsafePrivilege(fs api.FileSystem) {
	unsafepath := []string{"/etc/passwd", "/etc/crontab"}
	for i := 0; i < len(unsafepath); i++ {
		content, err := fs.Stat(unsafepath[i])
		if err == nil {
			priv := content.Mode()
			if FileWriteCheck(priv.String()) {
				AddResult(unsafepath[i], "UnsafePrivilege")
			}
		}
	}

}

//----------------------------容器逃逸相关------------------------------------

// 特权模式检查 --privileged --cap-add sys-admin 通过/proc/self/status中的CapEff判断，CapEff需要进行解析
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
func PrivilegeModeCheck(fs1 api.Container) {
	//--privileged模式下CapEff后几位全为f，其他要通过Parse解析
}

// 容器内挂载文件检查

func SensitiveFileMountCheck(fs1 api.Container) {
	spec, err := fs1.OCISpec()
	if err == nil {
		for _, mount := range spec.Mounts {
			if mount.Source == "/proc" { //procfs
				AddResult(mount.Destination, "UnSafeMount:Procfs(/proc/sys/kernel/core_pattern)")
			} else if strings.Contains(mount.Source, "lxcfs") {
				AddResult(mount.Destination, "UnSafeMount:lxcfs("+mount.Destination+")")
			}
		}
	} else {
		log.Info(err)
	}
}

// 检测docker remote api是否开启
func RemoteApiCheck(fs1 api.Container) {

}

// ----------------------------处理-------------------------------------------

func AddResult(s string, UnsafeType string) {
	result := &models.EscalateResult{
		ResultDetails: s,
		ResultType:    UnsafeType,
	}
	escalateLock.Lock()
	res = append(res, result)
	escalateLock.Unlock()
}

func GenerateImageRoport(image api.Image) error {
	if len(res) > 0 {
		detail, err := json.Marshal(res)
		if err == nil {
			Reportevent := report.ReportEvent{
				ID:         image.ID(),
				Time:       time.Now(),
				Level:      report.High,
				DetectType: report.Image,
				EventType:  report.Risk,
				AlertType:  report.Weakpass,
				GeneralDetails: []report.GeneralDetail{
					detail,
				},
			}
			err := report.DefaultReportClient().Report(Reportevent)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
func ImagesScanRun(fs api.Image) {

	go FindSuid(fs)
	go CheckEmptyPasswdRoot(fs)
	go UnsafePrivilege(fs)

	GenerateImageRoport(fs)
}

func GenerateContainerRoport(image api.Container) error {
	if len(res) > 0 {
		detail, err := json.Marshal(res)
		if err == nil {
			Reportevent := report.ReportEvent{
				ID:         image.ID(),
				Time:       time.Now(),
				Level:      report.High,
				DetectType: report.Image,
				EventType:  report.Risk,
				AlertType:  report.Weakpass,
				GeneralDetails: []report.GeneralDetail{
					detail,
				},
			}
			err := report.DefaultReportClient().Report(Reportevent)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
func ContainersScanRun(fs api.Container) {

	FindSuid(fs)
	CheckEmptyPasswdRoot(fs)
	UnsafePrivilege(fs)
	//PrivilegeModeCheck(fs)
	SensitiveFileMountCheck(fs)

	GenerateContainerRoport(fs)

}
