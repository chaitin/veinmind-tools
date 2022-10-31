module github.com/chaitin/veinmind-tools/plugins/go/veinmind-example

go 1.16

require (
	github.com/chaitin/libveinmind v1.3.1
	github.com/chaitin/veinmind-common-go v1.1.9
)

replace google.golang.org/grpc/naming => github.com/xiegeo/grpc-naming v1.29.1-alpha
