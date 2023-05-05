package detect

import (
	"context"
	"strings"
	"testing"

	"github.com/magiconair/properties/assert"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-webshell/pkg/filter"
)

func TestKit_Detect(t *testing.T) {
	detectKit, err := NewKit(context.Background(), WithToken("4689fb9bd16c13ab2fb80eeeb995ef6b"), WithDefaultClient())
	if err != nil {
		t.Error(err)
	}

	res, err := detectKit.Detect(FileInfo{
		Path:       "/tmp/1.php",
		Reader:     strings.NewReader("<?php eval($_GET[1]); ?>"),
		ScriptType: filter.PHP_TYPE,
	})
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, res.Data.RiskLevel, 20)
}
