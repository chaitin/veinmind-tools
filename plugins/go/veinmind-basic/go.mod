module github.com/chaitin/veinmind-tools/plugins/go/veinmind-basic

go 1.17

require (
	github.com/chaitin/libveinmind v1.2.1
	github.com/chaitin/veinmind-common-go v1.1.8-r0
	github.com/distribution/distribution v2.8.1+incompatible
	github.com/pkg/errors v0.9.1
	github.com/syndtr/gocapability v0.0.0-20200815063812-42c35b437635
)

require (
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.3-0.20220114050600-8b9d41f48198 // indirect
	github.com/opencontainers/runtime-spec v1.0.3-0.20210326190908-1c3f411f0417 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/cobra v1.4.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sync v0.0.0-20220513210516-0976fa681c29 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/xerrors v0.0.0-20220517211312-f3a8303e98df // indirect
)

replace google.golang.org/grpc/naming => github.com/xiegeo/grpc-naming v1.29.1-alpha
