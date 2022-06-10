package conf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultClient(t *testing.T) {
	c := DefaultConfClient()
	_, err := c.Pull(Sensitive)
	assert.Error(t, err)
}

func TestNewConfService(t *testing.T) {
	s, err := NewConfService()
	if err != nil {
		t.Error(err)
	}

	s.Store(Sensitive, []byte{0x01})
	b, err := s.Pull(Sensitive)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, b, []byte{0x01})
}
