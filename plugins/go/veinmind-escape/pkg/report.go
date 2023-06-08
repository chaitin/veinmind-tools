package pkg

import (
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/veinmind-common-go/service/report/event"
)

type CheckFunc func(api.FileSystem) ([]*event.EscapeDetail, error)

var (
	ImageCheckList     = make([]CheckFunc, 0)
	ContainerCheckList = make([]CheckFunc, 0)
)

const (
	WRITE             checkMode = 2
	READ              checkMode = 4
	KERNELPATTERN     string    = `([0-9]{1,})\.([0-9]{1,})\.([0-9]{1,})-[0-9]{1,}-[a-zA-Z]{1,}`
	SUDOREGEX         string    = `(\w{1,})\s\w{1,}=\(.*\)\s(.*)`
	CVEREASON         string    = "Your system has an insecure kernel version that is affected by a CVE vulnerability:"
	DOCKERAPIREASON   string    = "Docker remote API is opened which is can be used for escalating"
	SUDOREASON        string    = "This file is granted sudo privileges and can be used for escalating,you can check it in /etc/sudoers"
	MOUNTREASON       string    = "There are some sensitive files or directory mounted"
	READREASON        string    = "This file is sensitive and is readable to all users"
	WRITEREASON       string    = "This file is sensitive and is writable to all users"
	SUIDREASON        string    = "This file is granted suid privileges and belongs to root. And this file can be interacted with, there is a risk of elevation"
	EMPTYPASSWDREASON string    = "This user is privileged but does not have a password set"
	CAPREASON         string    = "There are unsafe linux capability granted"
)
