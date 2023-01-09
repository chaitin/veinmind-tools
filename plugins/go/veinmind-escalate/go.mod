module github.com/chaitin/veinmind-tools/plugins/go/veinmind-escalate

go 1.16

require (
	github.com/chaitin/libveinmind v1.5.1
	github.com/chaitin/veinmind-common-go v1.1.9
)

require (
	github.com/docker/docker v20.10.17+incompatible
	github.com/prometheus/procfs v0.7.3
	github.com/spf13/cobra v1.5.0 // indirect
)

replace google.golang.org/grpc/naming => github.com/xiegeo/grpc-naming v1.29.1-alpha
