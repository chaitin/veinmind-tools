package scanner

type ResultCode string

type Result struct {
	File    string `json:"file"`
	Version string `json:"version"`
}
