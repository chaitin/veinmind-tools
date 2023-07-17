package kernel

import (
	"bufio"
	"sort"
	"strconv"
	"strings"

	api "github.com/chaitin/libveinmind/go"
)

type SyscallEntry struct {
	Name string
	Addr uint64
}

type Ksyscall struct {
	SyscallList []*SyscallEntry
	SyscallMap  map[string]int
}

type KallsymsEntry struct {
	Addr uint64
	Type string
}

type KallSyms struct {
	Version         *KernelVersion
	SyscallEntry    *Ksyscall
	KallsymsMap     map[string]KallsymsEntry
	KernelTextRange uint64
}

func (syscallEntry *Ksyscall) Init() {
	syscallEntry.SyscallList = KcallList
	syscallEntry.SyscallMap = make(map[string]int)

	for i, entry := range syscallEntry.SyscallList {
		syscallEntry.SyscallMap[entry.Name] = i
	}
}

func (syscallEntry *Ksyscall) UpdateSyscall(name, typ string, addr uint64) bool {
	if typ == "T" || typ == "W" {
		index, exists := syscallEntry.SyscallMap[name]
		if !exists {
			index, exists = syscallEntry.SyscallMap[strings.TrimPrefix(name, "__x64_")]
		}
		if exists {
			syscallEntry.SyscallList[index].Addr = addr
		}
		return exists
	}
	return false
}

func (kallsyms *KallSyms) Init(apiFileSystem api.FileSystem) error {
	sort.Strings(SymsBuiltinList)
	syscallEntry := &Ksyscall{}
	syscallEntry.Init()
	kallsymsMap := make(map[string]KallsymsEntry)

	file, err := apiFileSystem.Open(KallsymsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	zeroAddressCount := 0
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		addr, err := strconv.ParseUint(fields[0], 16, 64)
		if err != nil {
			return err
		}

		if addr == 0 {
			zeroAddressCount++
			if zeroAddressCount > MaxZeroAddresses {
				return ErrKallsymsAddr
			}
		}

		if len(fields) > 3 && !strings.EqualFold(fields[3], "[ext4]") || len(fields) < 3 {
			continue
		}

		isUpdated := false
		if fields[1] < "a" {
			isUpdated = syscallEntry.UpdateSyscall(fields[2], fields[1], addr)
		}

		if !isUpdated && sort.SearchStrings(SymsBuiltinList, fields[2]) < len(SymsBuiltinList) {
			kallsymsMap[fields[2]] = KallsymsEntry{Addr: addr, Type: fields[1]}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	kallsyms.SyscallEntry = syscallEntry
	kallsyms.KallsymsMap = kallsymsMap

	return nil
}
