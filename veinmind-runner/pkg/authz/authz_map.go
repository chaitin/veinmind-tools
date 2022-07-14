package authz

import (
	"strings"
	"sync"
)

type policyMap struct {
	mtx      sync.Mutex
	policies map[string]Policy
}

func newPolicyMap() *policyMap {
	return &policyMap{
		policies: make(map[string]Policy),
	}
}

func (my *policyMap) Store(policy Policy) {
	my.mtx.Lock()
	my.policies[policy.Action] = policy
	my.mtx.Unlock()
}

func (my *policyMap) Load(action string) (Policy, bool) {
	my.mtx.Lock()
	policy, ok := my.policies[action]
	my.mtx.Unlock()

	return policy, ok
}

type handleMap struct {
	mtx       sync.Mutex
	handleIds map[string]struct{}
}

func newHandleMap() *handleMap {
	return &handleMap{
		handleIds: make(map[string]struct{}),
	}
}

func (my *handleMap) Store(action string) {
	my.mtx.Lock()
	my.handleIds[action] = struct{}{}
	my.mtx.Unlock()
}

func (my *handleMap) Delete(action string) {
	my.mtx.Lock()
	delete(my.handleIds, action)
	my.mtx.Unlock()
}

func (my *handleMap) Count(pattern string) int {
	my.mtx.Lock()
	count := 0
	for handleId, _ := range my.handleIds {
		if strings.HasPrefix(handleId, pattern) {
			count += 1
		}
	}
	my.mtx.Unlock()

	return count
}
