package pkg

import (
	"bufio"
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"regexp"
	"strconv"
)

type CVECheckFunc func(string, int, int, int)

var CVECheckList = []CVECheckFunc{
	checkCVECVE_2016_5195_DirtyCow,
	checkCVE_2020_14386,
	checkCVE2022_0847,
	checkCVE_2017_1000112,
	checkCVE_2021_22555,
}

func getVersion(fs api.FileSystem) ([]string, error) {
	content, err := fs.Open("/proc/version")
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer FileClose(content, err)
	scanner := bufio.NewScanner(content)
	for scanner.Scan() {
		complie := regexp.MustCompile(KERNELPATTERN)
		res1 := complie.FindStringSubmatch(scanner.Text())
		if len(res1) == 4 {
			return res1, nil
		}
	}
	return nil, nil
}

func ContainerCVECheck(fs api.FileSystem) error {
	version, err := getVersion(fs)
	if err != nil {
		log.Error(err)
		return err
	}
	KernelVersion, err := strconv.Atoi(version[0])
	if err != nil {
		return err
	}
	MajorRevision, err := strconv.Atoi(version[2])
	if err != nil {
		return err
	}
	MinorRevision, err := strconv.Atoi(version[3])
	if err != nil {
		return err
	}
	var versionString string
	for _, value := range version {
		versionString += value
	}
	for _, opts := range CVECheckList {
		opts(versionString, KernelVersion, MajorRevision, MinorRevision)
	}
	return nil
}

// 2.6.22 <= ver <= 4.8.3
func checkCVECVE_2016_5195_DirtyCow(version string, KernelVersion int, MajorRevision int, MinorRevision int) {
	if (KernelVersion == 2 && MajorRevision == 6 && MinorRevision >= 22) ||
		(KernelVersion == 2 && MajorRevision >= 7) ||
		(KernelVersion == 3) ||
		(KernelVersion == 4 && MajorRevision < 8) ||
		(KernelVersion == 4 && MajorRevision == 8 && MinorRevision <= 3) {
		AddResult("", CVEREASON+"DirtyCow", "UnsafeKernelVersion "+version)
	}
}

// 4.6 <= ver < 5.9
func checkCVE_2020_14386(version string, KernelVersion int, MajorRevision int, MinorRevision int) {
	if (KernelVersion == 4 && MajorRevision >= 6) ||
		(KernelVersion == 5 && MajorRevision < 9) {
		AddResult("", CVEREASON+"CVE-2020-14386", "UnsafeKernelVersion "+version)
	}
}

// 5.8 <= ver < 5.10.102 < ver < 5.15.25 <  ver <  5.16.11
func checkCVE2022_0847(version string, KernelVersion int, MajorRevision int, MinorRevision int) {
	if KernelVersion == 5 {
		if (MajorRevision >= 8 && MinorRevision < 10) ||
			(MajorRevision == 10 && MinorRevision < 102) ||
			(MajorRevision == 10 && MinorRevision > 102) ||
			(MajorRevision > 10 && MajorRevision < 15) ||
			(MajorRevision == 15 && MinorRevision < 25) ||
			(MajorRevision == 15 && MinorRevision > 25) ||
			(MajorRevision == 16 && MinorRevision < 11) {
			AddResult("", CVEREASON+"CVE-2023-0847", "UnsafeKernelVersion "+version)
		}
	}
}

// 4.4 <= ver<=4.13
func checkCVE_2017_1000112(version string, KernelVersion int, MajorRevision int, MinorRevision int) {
	if KernelVersion == 4 && MajorRevision >= 4 && MajorRevision <= 13 {
		AddResult("", CVEREASON+"CVE-2017-1000112", "UnsafeKernelVersion "+version)
	}
}

// 2.6.19 <= ver <= 5.12
func checkCVE_2021_22555(version string, KernelVersion int, MajorRevision int, MinorRevision int) {
	if (KernelVersion == 2 && MajorRevision == 6 && MinorRevision >= 19) ||
		(KernelVersion == 2 && MajorRevision >= 7) ||
		(KernelVersion == 3) ||
		(KernelVersion == 4) ||
		(KernelVersion == 5 && MajorRevision < 12) {
		AddResult("", CVEREASON+"CVE-2021-22555", "UnsafeKernelVersion "+version)
	}
}
