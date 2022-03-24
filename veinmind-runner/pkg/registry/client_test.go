package registry

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func TestList(t *testing.T) {
	c, err := NewRegistryClient("index.docker.io", nil)
	if err != nil {
		panic(err)
	}
	log.Println(c.GetRepos())
	d, err := c.GetRepo("ubuntu")
	m, _ := d.RawManifest()
	log.Println(string(m), err)

	r, err := c.Pull("ubuntu")
	if err != nil {
		panic(err)
	}

	resp, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(string(resp))
	}
}
