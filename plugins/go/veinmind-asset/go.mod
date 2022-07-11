module github.com/chaitin/veinmind-tools/plugins/go/veinmind-asset

go 1.18

require (
	github.com/aquasecurity/go-dep-parser v0.0.0-20220626060741-179d0b167e5f
	github.com/aquasecurity/trivy v0.29.2
	github.com/chaitin/libveinmind v1.1.2
	github.com/chaitin/veinmind-common-go v1.0.4
	github.com/spf13/cobra v1.5.0
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f
)

require (
	cloud.google.com/go/iam v0.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/go-containerregistry v0.10.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.3-0.20220303224323-02efb9a75ee1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/samber/lo v1.21.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.4.0 // indirect
	github.com/stretchr/testify v1.7.3 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/exp v0.0.0-20220407100705-7b9b53b0aca4 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/xerrors v0.0.0-20220517211312-f3a8303e98df // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	// containerd main
	github.com/containerd/containerd => github.com/containerd/containerd v1.6.1-0.20220606171923-c1bcabb45419
	// See https://github.com/moby/moby/issues/42939#issuecomment-1114255529
	github.com/docker/docker => github.com/docker/docker v20.10.3-0.20220224222438-c78f6963a1c0+incompatible
	google.golang.org/grpc/naming => github.com/xiegeo/grpc-naming v1.29.1-alpha
)
