package capability

import (
	"fmt"
	"testing"

	"github.com/chaitin/libveinmind/go/docker"
)

func TestIsPrivileged(t *testing.T) {
	d, _ := docker.New()
	ids, _ := d.ListContainerIDs()
	for _, id := range ids {
		c, _ := d.OpenContainerByID(id)
		fmt.Println(IsPrivileged(c))
	}
}
