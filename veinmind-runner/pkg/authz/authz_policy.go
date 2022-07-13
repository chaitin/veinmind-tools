package authz

import "sync"

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
