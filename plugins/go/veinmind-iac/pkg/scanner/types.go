package scanner

type Result struct {
	Risks []Risk `json:"risks"`
	*Rule `json:"rule"`
}

type Risk struct {
	StartLine int64
	EndLine   int64
	FilePath  string
	Original  string
}

type Rule struct {
	Id          string
	Name        string
	Description string
	Reference   string
	Severity    string
	Solution    string
	Type        string
}
