module github.com/chaitin/veinmind-tools/plugins/go/veinmind-unsafe-mount

go 1.16

require (
	github.com/chaitin/libveinmind v1.5.1
	github.com/chaitin/veinmind-common-go v1.3.0
)

require github.com/spf13/cobra v1.5.0 // indirect

replace google.golang.org/grpc/naming => github.com/xiegeo/grpc-naming v1.29.1-alpha
