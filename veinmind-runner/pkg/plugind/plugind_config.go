package plugind

type Config struct {
	Plugin []PluginConf `json:"plugin" ,toml:"plugin"`
}

type PluginConf struct {
	Name    string         `json:"name" ,toml:"name"`
	Service []*ServiceConf `json:"service" ,toml:"service"`
}

type ServiceConf struct {
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
