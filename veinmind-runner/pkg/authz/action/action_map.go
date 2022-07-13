package action

import (
	"sync"
)

type ActionMap struct {
	mtx     sync.Mutex
	actions map[string]struct{}
}

func NewMap() *ActionMap {
	return &ActionMap{
		actions: make(map[string]struct{}),
	}
}

func (my *ActionMap) Store(actionId string) {
	my.mtx.Lock()
	my.actions[actionId] = struct{}{}
	my.mtx.Unlock()
}

func (my *ActionMap) Delete(actionId string) {
	my.mtx.Lock()
	delete(my.actions, actionId)
	my.mtx.Unlock()
}

func (my *ActionMap) Count(pattern string) int {
	my.mtx.Lock()
	count := len(my.actions)
	my.mtx.Unlock()

	return count
}
