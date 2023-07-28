package kernel

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	api "github.com/chaitin/libveinmind/go"
)

type KernelVersion struct {
	Major int
	Minor int
	Patch int
}

func (version *KernelVersion) ParseVersionString(versionString string) error {
	versionParts := strings.Split(versionString, ".")
	if len(versionParts) < 3 {
		return ErrVersion
	}

	versionPartsInt := make([]int, len(versionParts))
	for i, v := range versionParts {
		ver, err := strconv.Atoi(v)
		if err != nil {
			return ErrVersion
		}
		versionPartsInt[i] = ver
	}

	version.Major = versionPartsInt[0]
	version.Minor = versionPartsInt[1]
	version.Patch = versionPartsInt[2]

	return nil
}

func (version *KernelVersion) GetKernelVersion(apiFileSystem api.FileSystem) error {
	file, err := apiFileSystem.Open(VersionPath)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	regex := regexp.MustCompile(`\d+\.\d+\.\d+`)
	versionString := regex.FindString(string(data))
	if versionString == "" {
		return ErrVersion
	}

	return version.ParseVersionString(versionString)
}
