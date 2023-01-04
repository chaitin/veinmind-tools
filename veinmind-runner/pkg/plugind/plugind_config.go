package plugind

type Manager struct {
	Plugins []Plugin `toml:"plugin"`
}

type Plugin struct {
	Name    string     `toml:"name"`
	Service []*Service `toml:"service"`
}

type Service struct {
	Name    string          `toml:"name"`
	Command string          `toml:"command"`
	Stdout  string          `toml:"stdout"`
	Stderr  string          `toml:"stderr"`
	Check   []*ServiceCheck `toml:"check"`
	Timeout int             `toml:"timeout"`
	Running bool
}

type ServiceCheck struct {
	Type  string `toml:"type"`
	Value string `toml:"value"`
}
