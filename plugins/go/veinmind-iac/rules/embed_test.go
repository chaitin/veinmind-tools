package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_EmbedBmRules(t *testing.T) {
	entries, err := RegoFile.ReadDir(".")
	require.NoError(t, err)
	assert.Greater(t, len(entries), 0)
}
