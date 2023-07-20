package service

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report/event"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-backdoor/kernel"
)

func removeDuplicates(details []*event.BackdoorDetail) []*event.BackdoorDetail {
	seen := make(map[string]bool, len(details))
	result := []*event.BackdoorDetail{}

	for _, entry := range details {
		if _, ok := seen[entry.Content]; !ok {
			seen[entry.Content] = true
			result = append(result, entry)
		}
	}

	return result
}

func findModule(
	addr uint64,
	kcore *kernel.KcoreMemory,
	kmod *kernel.KernelModules,
	version *kernel.KernelVersion,
) (kernel.ModuleInfo, error) {
	pos := kmod.BinarySearch(addr)
	if pos < 0 {
		return kernel.ModuleInfo{}, kernel.ErrInvalidNum
	} else if pos != 0 {
		prevModule := *kmod.ModuleList[pos-1]
		moduleRange := prevModule.Addr + prevModule.Size
		if addr < moduleRange {
			return prevModule, nil
		} else if pos != len(kmod.ModuleList) && addr == kmod.ModuleList[pos].Addr {
			return *kmod.ModuleList[pos], nil
		}
	}

	modInfo, err := kcore.FindModule(addr, kmod.ModOffset, version)
	if err != nil {
		return kernel.ModuleInfo{}, err
	}
	kmod.Insert(pos, &modInfo)

	return modInfo, nil
}

func rootkitPathCheck(apiFileSystem api.FileSystem, name, path string) *event.BackdoorDetail {
	rootkitFileInfo, err := apiFileSystem.Lstat(path)
	if err != nil {
		return nil
	}

	fileDetail, err := file2FileDetail(rootkitFileInfo, path)
	if err != nil {
		return nil
	}

	res := &event.BackdoorDetail{
		FileDetail:  fileDetail,
		Content:     name + ": " + path,
		Description: kernel.DefaultDescription,
	}

	return res
}

func rootkitRuleCheck(apiFileSystem api.FileSystem) (bool, []*event.BackdoorDetail) {
	var res []*event.BackdoorDetail
	check := false

	for _, rootkitInfo := range kernel.RootkitRules {
		checkPaths := append(rootkitInfo.File, rootkitInfo.Dir...)

		for _, path := range checkPaths {
			if checkRes := rootkitPathCheck(
				apiFileSystem,
				rootkitInfo.Name,
				path,
			); checkRes != nil {
				check = true
				res = append(res, checkRes)
			}
		}
	}

	return check, res
}

func rootkitLKMCheck(apiFileSystem api.FileSystem) (bool, []*event.BackdoorDetail) {
	check := false
	var res []*event.BackdoorDetail

	dirLKM, err := apiFileSystem.Lstat(kernel.LKMDir)
	if err != nil || !dirLKM.IsDir() {
		return false, nil
	}

	apiFileSystem.Walk(kernel.LKMDir, func(path string, info fs.FileInfo, err error) error {
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".so" || ext == ".ko" || ext == ".ko.xz" {
			for _, lkm := range kernel.BadLKM {
				if lkm == strings.TrimSuffix(strings.ToLower(filepath.Base(path)), ext) {
					check = true
					fileDetail, err := file2FileDetail(info, path)
					if err != nil {
						return err
					}
					res = append(res, &event.BackdoorDetail{
						FileDetail:  fileDetail,
						Content:     path,
						Description: kernel.DefaultDescription,
					})
				}
			}
		}
		return err
	})

	return check, res
}

func initEnv(
	apiFileSystem api.FileSystem,
	kallsyms *kernel.KallSyms,
	kcore *kernel.KcoreMemory,
	kmod *kernel.KernelModules,
	version *kernel.KernelVersion,
) error {
	err := kallsyms.Init(apiFileSystem)
	if err != nil {
		return err
	}

	tmpText, ok := kallsyms.KallsymsMap["_text"]
	if !ok {
		return kernel.ErrKallsymsInit
	}
	kcore.TextAddr = tmpText.Addr
	err = kcore.Init(apiFileSystem, kcore.TextAddr)
	if err != nil {
		return err
	}
	if kcore.TextSize != 0 {
		kcore.KernelTextRange = kcore.TextAddr + kcore.TextSize
	}

	err = kmod.Init(apiFileSystem)
	if err != nil {
		return nil
	}
	info, err := apiFileSystem.Stat(kernel.KernelModulesPath)
	if err == nil {
		kmod.ModDetail, _ = file2FileDetail(info, kernel.KernelModulesPath)
	}

	err = version.GetKernelVersion(apiFileSystem)
	if err != nil {
		return err
	}
	fmt.Println(version)

	return nil
}

// ref: https://github.com/grayddq/GScan/blob/master/lib/plugins/Rootkit_Analysis.py
func rootkitBackdoorCheck(apiFileSystem api.FileSystem) (bool, []*event.BackdoorDetail) {
	var res []*event.BackdoorDetail
	check := false

	tmpCheck, tmpRes := rootkitRuleCheck(apiFileSystem)
	check = check || tmpCheck
	res = append(res, tmpRes...)

	tmpCheck, tmpRes = rootkitLKMCheck(apiFileSystem)
	check = check || tmpCheck
	res = append(res, tmpRes...)

	res = removeDuplicates(res)

	return check, res
}

func rootkitContainerCheck(apiFileSystem api.FileSystem) (bool, []*event.BackdoorDetail) {
	var res []*event.BackdoorDetail
	check := false

	tmpCheck, tmpRes := rootkitBackdoorCheck(apiFileSystem)
	check = check || tmpCheck
	res = append(res, tmpRes...)

	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		return check, res
	}

	kallsyms := &kernel.KallSyms{}
	kcore := &kernel.KcoreMemory{}
	kmod := &kernel.KernelModules{}
	version := &kernel.KernelVersion{}
	err := initEnv(
		apiFileSystem,
		kallsyms,
		kcore,
		kmod,
		version,
	)
	defer kcore.FileHandle.Close()
	if err != nil {
		return check, res
	}

	syscallTable, ok := kallsyms.KallsymsMap["sys_call_table"]
	if !ok {
		return check, res
	}

	memData, err := kcore.Read(syscallTable.Addr, uint64(len(kallsyms.SyscallEntry.SyscallList)*8+8))
	if err != nil {
		return check, res
	}

	for sysNum, info := range kallsyms.SyscallEntry.SyscallList {
		addr := uint64(kernel.BytesToUint(memData[sysNum*8 : sysNum*8+8]))
		if addr < uint64(0xffff800000000000) || addr > uint64(0xffffffffffff8000) {
			continue
		}

		if addr != info.Addr {
			if kcore.TextAddr <= addr && addr < kcore.KernelTextRange {
				continue
			}
			module, err := findModule(addr, kcore, kmod, version)
			if err != nil {
				continue
			}
			if module.Name != "" {
				check = true
				res = append(res, &event.BackdoorDetail{
					FileDetail:  kmod.ModDetail,
					Content:     module.Name + ": " + strconv.FormatUint(module.Addr, 16),
					Description: kernel.DefaultDescription,
				})
			}
		}

	}

	res = removeDuplicates(res)

	return check, res
}

func init() {
	ImageCheckFuncMap["rootkit"] = rootkitBackdoorCheck
	ContainerCheckFuncMap["rootkit"] = rootkitContainerCheck
}
