module github.com/chaitin/veinmind-tools/plugins/go/veinmind-escalate

go 1.16

require github.com/chaitin/libveinmind v1.3.1

require (
	github.com/chaitin/veinmind-common-go v1.1.9
	github.com/chaitin/veinmind-tools/plugins/go/veinmind-unsafe-mount v0.0.0-20221112032047-46463828bdeb
	github.com/docker/docker v20.10.17+incompatible
	github.com/opencontainers/runtime-spec v1.0.3-0.20210326190908-1c3f411f0417
)

replace google.golang.org/grpc/naming => github.com/xiegeo/grpc-naming v1.29.1-alpha
