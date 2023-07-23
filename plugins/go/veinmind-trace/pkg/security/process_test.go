package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHideMount(t *testing.T) {
	content1 := `sysfs /sys sysfs rw,nosuid,nodev,noexec,relatime 0 0
proc /proc proc rw,nosuid,nodev,noexec,relatime 0 0`
	content2 := `sysfs /sys sysfs rw,nosuid,nodev,noexec,relatime 0 0
proc /proc/1 proc rw,nosuid,nodev,noexec,relatime 0 0`
	a, _ := hasMount(content1)
	b, _ := hasMount(content2)
	assert.Equal(t, false, a)
	assert.Equal(t, true, b)
}

func TestHasPtraceProcess(t *testing.T) {
	content1 := `Name:   rsyslogd
Umask:  0022
State:  S (sleeping)
Tgid:   3330111
Ngid:   0
Pid:    3330111
PPid:   1
TracerPid:      0
`
	content2 := `Name:   rsyslogd
Umask:  0022
State:  S (sleeping)
Tgid:   3330111
Ngid:   0
Pid:    3330111
PPid:   1
TracerPid:    2
`
	assert.Equal(t, false, HasPtraceProcess(content1))
	assert.Equal(t, true, HasPtraceProcess(content2))
}
