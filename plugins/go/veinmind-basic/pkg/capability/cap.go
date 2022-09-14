package capability

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/plugin/log"
)

func IsPrivileged(container api.Container) bool {
	state, err := container.OCIState()
	if err != nil {
		log.Error(err)
		return false
	}

	if state.Pid == 0 {
		return false
	}

	status, err := ioutil.ReadFile(filepath.Join(func() string {
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
