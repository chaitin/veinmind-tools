package log

import (
	"github.com/sirupsen/logrus"
)

var (
	modules map[string]*Module
)

type Module struct {
	name string
	*logrus.Entry
}

type ModuleKey string

const (
	CmdModuleKey       ModuleKey = "cmd"
	AIAnalyzerKey      ModuleKey = "ai-analyzer"
	ScanModuleKey      ModuleKey = "scan"
	AuthzModuleKey     ModuleKey = "authz"
	ContainerModuleKey ModuleKey = "container"
	GitModuleKey       ModuleKey = "git"
	PlugindModuleKey   ModuleKey = "plugind"
	RegistryModuleKey  ModuleKey = "registry"
	ReporterModuleKey  ModuleKey = "reporter"
	TargetModuleKey    ModuleKey = "target"
)

func (m ModuleKey) String() string {
	return string(m)
}

func (m Module) String() string {
	return m.name
}

func RegisterModule(name string) error {
	logger := logrus.New()

	// set formatter
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		ForceQuote:      true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	// set level
	logger.SetLevel(logrus.InfoLevel)

	// entry
	entry := logger.WithFields(map[string]interface{}{
		"module": name,
	})

	modules[name] = &Module{
		name:  name,
		Entry: entry,
	}

	return nil
}

func init() {
	// register
	modules = make(map[string]*Module)
	_ = RegisterModule(CmdModuleKey.String())
	_ = RegisterModule(AIAnalyzerKey.String())
	_ = RegisterModule(ScanModuleKey.String())
	_ = RegisterModule(AuthzModuleKey.String())
	_ = RegisterModule(ContainerModuleKey.String())
	_ = RegisterModule(GitModuleKey.String())
	_ = RegisterModule(PlugindModuleKey.String())
	_ = RegisterModule(RegistryModuleKey.String())
	_ = RegisterModule(ReporterModuleKey.String())
	_ = RegisterModule(TargetModuleKey.String())
}

func GetModule(key ModuleKey) *Module {
	if m, ok := modules[key.String()]; ok {
		return m
	} else {
		return defaultModule
	}
}
