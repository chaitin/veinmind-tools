package report

type Event struct {
	Source      string `json:"source" yaml:"source"`
	Destination string `json:"destination" yaml:"destination"`
	Type        string `json:"type" yaml:"type"`
}
