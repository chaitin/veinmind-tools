package conf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultClient(t *testing.T)  {
	c := DefaultConfClient()
	_, err :=  c.Pull(Sensitive)
	assert.Error(t, err)
}


