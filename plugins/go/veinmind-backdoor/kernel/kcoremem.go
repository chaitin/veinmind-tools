package kernel

import (
	"bytes"
	"debug/elf"
	"encoding/binary"

	api "github.com/chaitin/libveinmind/go"
)

type KcoreMemory struct {
	FileHandle        api.File
	FileHeader        *elf.Header64
	FileToVaddrOffset int64
	TextAddr          uint64
	TextSize          uint64
	KernelTextRange   uint64
	readBuffer        []byte
}

func (kcore *KcoreMemory) vAddrSeek(offset uint64) error {
	realAddr := offset + uint64(kcore.FileToVaddrOffset)
	if realAddr > uint64(MaxInt64) {
		return ErrTooBig
	}
	_, err := kcore.FileHandle.Seek(int64(realAddr), 0)
	return err
}

func (kcore *KcoreMemory) Read(offset, size uint64) ([]byte, error) {
	if len(kcore.readBuffer) < int(size) {
		kcore.readBuffer = make([]byte, size)
	}
	data := kcore.readBuffer[:size]

	err := kcore.vAddrSeek(offset)
	if err != nil {
		return nil, err
	}
	n, err := kcore.FileHandle.Read(data)
	if err != nil {
		return nil, err
	}
	if n != int(size) {
		return nil, ErrSizeMissMatch
	}

	return data, nil
}

func (kcore *KcoreMemory) ReadI(offset uint64) (uint32, error) {
	data, err := kcore.Read(offset, 4)
	if err != nil {
		return 0, err
	}

	return uint32(BytesToUint(data)), nil
}

func (kcore *KcoreMemory) ReadQ(offset uint64) (uint64, error) {
	data, err := kcore.Read(offset, 8)
	if err != nil {
		return 0, err
	}

	return uint64(BytesToUint(data)), nil
}

func (kcore *KcoreMemory) FindModule(
	addr uint64,
	modOffset int,
	version *KernelVersion,
) (ModuleInfo, error) {
	andAddr := addr & uint64(0xffffffffffff8000)
	startAddr := andAddr - 0x4000
	memData, err := kcore.Read(startAddr, 0x20*0x1000+8)
	if err != nil {
		return ModuleInfo{}, err
	}

	for offset := 0; offset <= 0x20*0x1000; offset += 0x10 {
		entry := uint64(BytesToUint(memData[offset : offset+8]))
		if entry == 0 || entry%16 != 0 {
			continue
		}
		readAddr := entry + uint64(modOffset)
		dataQ, err := kcore.ReadQ(readAddr)
		if err != nil || dataQ == entry {
			continue
		}
		dataI, err := kcore.ReadI(entry)
		if err != nil || int32(dataI) > 3 || int32(dataI) < 0 {
			continue
		}

		var modSize uint32
		var modBase uint64
		var modName string
		if version.Major < 4 {
			readAddr := entry + 0x138
			modBase, err = kcore.ReadQ(readAddr)
			if err != nil || modBase != 0 {
				continue
			}
			readAddr = entry + 0x144
			modSize, err = kcore.ReadI(readAddr)
			if err != nil {
				continue
			}

			for modOffset := 0x100; modOffset <= 0x240; modOffset += 0x8 {
				readAddr := entry + uint64(modOffset)
				data, err := kcore.Read(readAddr, 48)
				if err != nil {
					continue
				}
				var qVar [3]uint64
				var iVar [6]uint32
				for i := 0; i < 3; i++ {
					qVar[i] = uint64(BytesToUint(data[i*8 : i*8+8]))
				}
				for i := 0; i < 6; i++ {
					iVar[i] = uint32(BytesToUint(data[i*4+24 : i*4+28]))
				}
				band := entry & 0xfffffffffff00000
				if (qVar[0]&0xfffffffffff00000) == band &&
					((qVar[1]&0xfffffffffff00000) == band || qVar[1] == 0) &&
					(qVar[2]&0xfffffffffff80000) == (entry&0xfffffffffff80000) {
					for i := 0; i < 6; i++ {
						if iVar[i] > 0x10000 {
							continue
						}
					}
					modBase = qVar[2]
					if iVar[1] > 0 {
						modSize = iVar[0]
					} else if iVar[4] > 0 {
						modSize = iVar[4]
					} else {
						modSize = 0x4000
					}
				}
			}
		} else if version.Major == 4 && version.Minor < 5 {
			for modOffset := 0x100; modOffset <= 0x240; modOffset += 0x40 {
				readAddr := entry + uint64(modOffset)
				data, err := kcore.Read(readAddr, 32)
				if err != nil {
					continue
				}
				var qVar [2]uint64
				var iVar [4]uint32
				for i := 0; i < 2; i++ {
					qVar[i] = uint64(BytesToUint(data[i*8 : i*8+8]))
				}
				for i := 0; i < 4; i++ {
					iVar[i] = uint32(BytesToUint(data[i*4+16 : i*4+20]))
				}
				if ((qVar[0]&0xfffffffffff00000) == (entry&0xfffffffffff00000) || qVar[0] == 0) &&
					(qVar[1]&0xfffffffffff80000) == (entry&0xfffffffffff80000) &&
					iVar[1] > 0 && iVar[3] > 0 {
					flag := true
					for i := 0; i < 4; i++ {
						if (iVar[i] & 0xfff) != 0 {
							flag = false
							break
						}
					}
					if flag {
						modBase = qVar[1]
						modSize = iVar[1]
					}
				}
			}
		}

		if modBase == 0 {
			for modOffset := 0x100; modOffset <= 0x240; modOffset += 0x40 {
				readAddr := entry + uint64(modOffset)
				data, err := kcore.Read(readAddr, 24)
				if err != nil {
					continue
				}
				var qVar uint64
				var iVar [4]uint32
				qVar = uint64(BytesToUint(data[0:8]))
				for i := 0; i < 4; i++ {
					iVar[i] = uint32(BytesToUint(data[i*4+8 : i*4+12]))
				}
				if (qVar&0xfffffffffff80000) == (entry&0xfffffffffff80000) &&
					(iVar[3]&0xfff) == 0 {
					flag := true
					for i := 0; i < 3; i++ {
						if (iVar[i] & 0xfff) != 0 {
							flag = false
							break
						}
					}
					if flag {
						modBase = qVar
						modSize = iVar[0]
					}
				}
			}
		}

		modRange := modBase + uint64(modSize)
		if (modBase != 0) && modBase > uint64(0xffff800000000000) &&
			modSize < 0x20000 && modSize > 0xa00 && modRange > modBase {
			readAddr := entry + 0x18
			data, err := kcore.Read(readAddr, 64)
			if err != nil {
				continue
			}
			startPos := bytes.IndexByte(data, 0)
			if startPos < 0 {
				continue
			}
			modName = string(data[:startPos])
			if modBase <= addr && addr < modRange {
				return ModuleInfo{
					Addr: modBase,
					Size: uint64(modSize),
					Name: modName,
				}, nil
			} else {
				if addr < modBase {
					break
				}
			}
		}
	}

	return ModuleInfo{}, ErrNoModule
}

func (kcore *KcoreMemory) Init(
	apiFileSystem api.FileSystem,
	textAddr uint64,
) error {
	file, err := apiFileSystem.Open(KcorePath)
	if err != nil {
		return err
	}

	elfFile, err := elf.NewFile(file)
	if err != nil {
		return err
	}
	defer elfFile.Close()

	header := &elf.Header64{}
	if err := binary.Read(file, binary.LittleEndian, header); err != nil {
		return err
	}
	offset, textSize, err := calOffsetAndSize(elfFile.Progs, textAddr)
	if err != nil {
		return err
	}

	kcore.FileHandle = file
	kcore.FileHeader = header
	kcore.FileToVaddrOffset = offset
	kcore.TextSize = textSize

	return nil
}

func calOffsetAndSize(progs []*elf.Prog, textAddr uint64) (int64, uint64, error) {
	var offset int64
	var textSize uint64

	for _, prog := range progs {
		if prog.Vaddr == textAddr {
			offset = int64(prog.Off) - int64(prog.Vaddr)
			textSize = prog.Memsz
			break
		} else if offset == 0 && prog.Vaddr > uint64(0xffffffff00000000) {
			offset = int64(prog.Off) - int64(prog.Vaddr)
		}
	}
	if offset == 0 {
		return 0, 0, ErrKcoreInit
	}

	return offset, textSize, nil
}
