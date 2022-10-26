package scanner

type ResultCode string

const (
	Vulnerable     ResultCode = "Vulnerable"
	FixedVersion   ResultCode = "FixedVersion"
	NotDetected    ResultCode = "NotDetected"
	JarDetectDepth            = 16
)

type Result struct {
	Code        ResultCode `json:"code"`
	File        string     `json:"file"`
	DisplayPath string     `json:"parent"`
}
