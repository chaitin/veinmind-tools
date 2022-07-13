package action

import (
	"sync"
)

type Map struct {
	mtx     sync.Mutex
	actions map[string]struct{}
}

func NewMap() *Map {
	return &Map{
		actions: make(map[string]struct{}),
	}
}

func (my *Map) Store(actionId string) {
	my.mtx.Lock()
	my.actions[actionId] = struct{}{}
	my.mtx.Unlock()
}

func (my *Map) Delete(actionId string) {
	my.mtx.Lock()
	delete(my.actions, actionId)
	my.mtx.Unlock()
}

func (my *Map) Count(pattern string) int {
	my.mtx.Lock()
	count := len(my.actions)
	my.mtx.Unlock()

	return count
}
