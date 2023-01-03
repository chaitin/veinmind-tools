package cache

import (
	"github.com/gogf/gf/container/gset"
)

var WhitePath = whitePathCache{
	mem: gset.NewStrSet(true),
}

type whitePathCache struct {
	mem *gset.StrSet
}

func (c *whitePathCache) Contains(key string) bool {
	return c.mem.Contains(key)
}

func (c *whitePathCache) Add(key string) {
	c.mem.Add(key)
}
