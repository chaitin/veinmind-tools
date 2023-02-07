package cache

import (
	"github.com/gogf/gf/os/gmutex"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-sensitive/rule"
)

var HashRule = hashRuleCache{
	mux: gmutex.New(),
	mem: make(map[string]map[int64]rule.Rule),
}

type hashRuleCache struct {
	mux *gmutex.Mutex
	mem map[string](map[int64]rule.Rule)
}

func (c *hashRuleCache) Get(key string) (map[int64]rule.Rule, bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	val, ok := c.mem[key]
	return val, ok
}

func (c *hashRuleCache) Set(key string, rules map[int64]rule.Rule) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if key == "" {
		return
	}

	c.mem[key] = rules
}

func (c *hashRuleCache) SetOrAppend(key string, r rule.Rule) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if key == "" {
		return
	}

	if _, ok := c.mem[key]; !ok {
		c.mem[key] = make(map[int64]rule.Rule)
	}
	c.mem[key][r.Id] = r
}
