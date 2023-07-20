package kernel

import (
	"testing"

	"gotest.tools/assert"
)

func TestBytesToUint(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected uint
	}{
		{
			name:     "One Byte",
			input:    []byte{0x1},
			expected: 1,
		},
		{
			name:     "Two Bytes",
			input:    []byte{0x0, 0x1},
			expected: 256,
		},
		{
			name:     "Four Bytes",
			input:    []byte{0x0, 0x0, 0x0, 0x1},
			expected: 16777216,
		},
		{
			name:     "Eight Bytes",
			input:    []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1},
			expected: 72057594037927936,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := BytesToUint(tc.input)
			if output != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, output)
			}
		})
	}
}

func TestParseVersionString(t *testing.T) {
	testData := []struct {
		versionString string
		expect        KernelVersion
		expectErr     error
	}{
		{
			versionString: "5.15.0",
			expect:        KernelVersion{Major: 5, Minor: 15, Patch: 0},
			expectErr:     nil,
		},
		{
			versionString: "5.15",
			expect:        KernelVersion{},
			expectErr:     ErrVersion,
		},
		{
			versionString: "5.15.t",
			expect:        KernelVersion{},
			expectErr:     ErrVersion,
		},
	}

	for _, data := range testData {
		var version KernelVersion
		err := version.ParseVersionString(data.versionString)

		assert.Equal(t, data.expectErr, err)
		if err == nil {
			assert.Equal(t, data.expect, version)
		}
	}
}

func TestBinarySearch(t *testing.T) {
	kmod := &KernelModules{
		ModuleList: []*ModuleInfo{
			{Addr: 1, Size: 100, Name: "Module1"},
			{Addr: 2, Size: 200, Name: "Module2"},
			{Addr: 3, Size: 300, Name: "Module3"},
		},
	}

	testCases := []struct {
		desc     string
		addr     uint64
		expected int
	}{
		{"addr exists", 2, 1},
		{"addr does not exist", 4, 3},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res := kmod.BinarySearch(tC.addr)
			if res != tC.expected {
				t.Fatalf("Expected %d, but got %d", tC.expected, res)
			}
		})
	}
}

func TestInsert(t *testing.T) {
	kmod := &KernelModules{
		ModuleList: []*ModuleInfo{
			{Addr: 1, Size: 100, Name: "Module1"},
			{Addr: 3, Size: 300, Name: "Module3"},
		},
	}

	module := &ModuleInfo{Addr: 2, Size: 200, Name: "Module2"}
	kmod.Insert(1, module)

	if kmod.ModuleList[1] != module {
		t.Fatalf("Expected %v, but got %v", module, kmod.ModuleList[1])
	}
}

func TestParseModuleInfo(t *testing.T) {
	testCases := []struct {
		desc     string
		fields   []string
		expected *ModuleInfo
		hasErr   bool
	}{
		{"valid fields", []string{"Module1", "100", "", "", "", "0x1"}, &ModuleInfo{1, 100, "Module1"}, false},
		{"invalid size", []string{"Module1", "abc", "", "", "", "0x1"}, nil, true},
		{"invalid addr", []string{"Module1", "100", "", "", "", "fgh"}, nil, true},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			name, size, addr, err := parseModuleInfo(tC.fields)
			if tC.hasErr && err == nil {
				t.Fatalf("Expected error, but got none")
			}

			if !tC.hasErr && err != nil {
				t.Fatalf("Did not expect error, but got %v", err)
			}

			if tC.expected != nil {
				if name != tC.expected.Name || size != tC.expected.Size || addr != tC.expected.Addr {
					t.Fatalf("Expected %v, but got %v %v %v", tC.expected, name, size, addr)
				}
			}
		})
	}
}

func TestUpdateSyscall(t *testing.T) {
	kcall := &Ksyscall{}
	kcall.Init()

	testCases := []struct {
		desc         string
		name         string
		typ          string
		addr         uint64
		expectResult bool
	}{
		{"valid syscall update", "sys_read", "T", 1234, true},
		{"invalid syscall type", "syscall_name", "invalid_type", 1234, false},
		{"non-existing syscall", "non_existing_syscall", "T", 1234, false},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			isSuccess := kcall.UpdateSyscall(tc.name, tc.typ, tc.addr)
			if isSuccess != tc.expectResult {
				t.Errorf("UpdateSyscall(%s, %s, %d) isSuccess = %v, wantResult %v", tc.name, tc.typ, tc.addr, isSuccess, tc.expectResult)
			}
		})
	}
}
