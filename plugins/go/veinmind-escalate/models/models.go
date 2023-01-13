package models

type EscalateResult struct {
	Target string `json:"target"`
	Reason string `json:"reason"`
	Detail string `json:"detail"`
}
