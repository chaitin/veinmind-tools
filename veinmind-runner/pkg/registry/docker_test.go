package registry

import (
	"log"
	"testing"
)

func TestList(t *testing.T) {
	c, err := NewRegistryDockerClient()
	if err != nil {
		panic(err)
	}
	switch v := c.(type) {
	case *RegistryDockerClient:
		log.Println(v.GetRepos("127.0.0.1:5000"))
		d, err := v.GetRepo("ubuntu")
		if err != nil {
			t.Error(err)
		}

		m, _ := d.RawManifest()
		log.Println(string(m), err)

		_, err = c.Pull("ubuntu")
		if err != nil {
			t.Error(err)
		}
	}

}
