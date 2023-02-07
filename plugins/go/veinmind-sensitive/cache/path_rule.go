package cache

import (
	"github.com/gogf/gf/os/gmutex"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-sensitive/rule"
)

var PathRule = pathRuleCache{
	mux: gmutex.New(),
	mem: make(map[string]map[int64]rule.Rule),
}

type pathRuleCache = hashRuleCache
