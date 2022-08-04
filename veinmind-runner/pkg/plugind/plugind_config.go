package plugind

type Manager struct {
	Plugins []Plugin `json:"plugin" ,toml:"plugin"`
}

type Plugin struct {
	Name    string     `json:"name" ,toml:"name"`
	Service []*Service `json:"service" ,toml:"service"`
}

type Service struct {
	Name    string          `json:"name" ,toml:"name"`
	Command string          `json:"command" ,toml:"command"`
	Stdout  string          `json:"stdout" ,toml:"stdout"`
	Stderr  string          `json:"stderr" ,toml:"stderr"`
	Check   []*ServiceCheck `json:"check" ,toml:"check"`
	Timeout int             `json:"timeout" ,toml:"timeout"`
}

type ServiceCheck struct {
	Type  string `json:"type" ,toml:"type"`
	Value string `json:"value" ,toml:"value"`
}
