package authz

import (
	"strings"
	"sync"
)

type actionMap struct {
	sync.Map
}

func (a *actionMap) Count(pattern string) int {
	var count int
	a.Range(func(key, value any) bool {
		imageActionId := key.(string)
		if strings.HasPrefix(imageActionId, pattern) {
			count += 1
		}

		return true
	})

	return count
}
