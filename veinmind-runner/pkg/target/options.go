package target

type Options struct {
	Thread                int
	TempPath              string
	ConfigPath            string
	ResourcePath          string
	ParallelContainerMode bool
	Insecure              bool
	SpecFlags             []string

	// For Iac Scan
	IacFileType      string
	IacLimitSize     int64
	IacSshPath       string
	IacProxy         string
	IacKubeConfig    string
	IacKubeNameSpace string

	// For Registry Scan
	CatalogFilterRegex string
}

type Option func(*Options)

func WithThread(thread int) Option {
	return func(o *Options) {
		o.Thread = thread
	}
}

func WithSpecFlags(flags []string) Option {
	return func(o *Options) {
		o.SpecFlags = flags
	}
}

func WithTempPath(path string) Option {
	return func(o *Options) {
		o.TempPath = path
	}
}

func WithConfigPath(path string) Option {
	return func(o *Options) {
		o.ConfigPath = path
	}
}

func WithFilterRegex(regex string) Option {
	return func(o *Options) {
		o.CatalogFilterRegex = regex
	}
}

func WithResourcePath(path string) Option {
	return func(o *Options) {
		o.ResourcePath = path
	}
}

func WithParallelContainerMode(mode bool) Option {
	return func(o *Options) {
		o.ParallelContainerMode = mode
	}
}

func WithInsecure(insecure bool) Option {
	return func(o *Options) {
		o.Insecure = insecure
	}
}

func WithIacFileType(fileType string) Option {
	return func(o *Options) {
		o.IacFileType = fileType
	}
}

func WithIacLimitSize(iacSize int64) Option {
	return func(o *Options) {
		o.IacLimitSize = iacSize
	}
}

func WithIacSshPath(path string) Option {
	return func(o *Options) {
		o.IacSshPath = path
	}
}

func WithIacProxy(proxy string) Option {
	return func(o *Options) {
		o.IacProxy = proxy
	}
}

func WithIacKubeConfig(configPath string) Option {
	return func(o *Options) {
		o.IacKubeConfig = configPath
	}
}

func WithIacKubeNameSpace(namespace string) Option {
	return func(o *Options) {
		o.IacKubeNameSpace = namespace
	}
}
