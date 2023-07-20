package kernel

import (
	"bufio"
	"fmt"
	"sort"
	"strconv"
	"strings"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report/event"
)

type ModuleInfo struct {
	Addr uint64
	Size uint64
	Name string
}

type KernelModules struct {
	ModuleList []*ModuleInfo
	ModDetail  event.FileDetail
	ModOffset  int
}

func (kmod *KernelModules) BinarySearch(addr uint64) int {
	left, right := 0, len(kmod.ModuleList)
	if kmod.ModuleList == nil {
		return -1
	}

	for left < right {
		mid := left + (right-left)/2
		if kmod.ModuleList[mid].Addr < addr {
			left = mid + 1
		} else {
			right = mid
		}
	}

	return left
}

func (kmod *KernelModules) Insert(index int, value *ModuleInfo) {
	kmod.ModuleList = append(kmod.ModuleList, nil)
	copy(kmod.ModuleList[index+1:], kmod.ModuleList[index:])
	kmod.ModuleList[index] = value
}

func (kmod *KernelModules) Init(apiFileSystem api.FileSystem) error {
	moduleList := []*ModuleInfo{}
	modDetail := event.FileDetail{}
	modOffset := ModuleOffset

	file, err := apiFileSystem.Open(KernelModulesPath)
	if err != nil {
		return err
	}
	defer file.Close()

	zeroAddressCount := 0
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		moduleName, moduleSize, moduleAddress, parseErr := parseModuleInfo(fields)
		if parseErr != nil {
			return parseErr
		}

		if moduleAddress == 0 {
			zeroAddressCount++
			if zeroAddressCount > MaxZeroAddresses {
				return ErrModulesAddr
			}
		}

		if moduleName != "" && moduleSize != 0 && moduleAddress != 0 {
			newModule := &ModuleInfo{
				Addr: moduleAddress,
				Size: moduleSize,
				Name: moduleName,
			}
			moduleList = append(moduleList, newModule)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	sort.Slice(moduleList, func(i, j int) bool {
		return moduleList[i].Addr < moduleList[j].Addr
	})

	kmod.ModuleList = moduleList
	kmod.ModDetail = modDetail
	kmod.ModOffset = modOffset

	return nil
}

func parseModuleInfo(fields []string) (name string, size uint64, addr uint64, err error) {
	size, err = strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return "", 0, 0, err
	}

	addrStr := fields[5]
	if !strings.HasPrefix(addrStr, "0x") {
		return "", 0, 0, fmt.Errorf("address should start with '0x'")
	}

	addr, err = strconv.ParseUint(strings.TrimPrefix(addrStr, "0x"), 16, 64)
	if err != nil {
		return "", 0, 0, err
	}

	name = fields[0]
	return name, size, addr, nil
}
