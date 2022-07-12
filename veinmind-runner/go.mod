module github.com/chaitin/veinmind-tools/veinmind-runner

go 1.16

require (
	github.com/chaitin/libveinmind v1.1.1
	github.com/chaitin/veinmind-common-go v1.0.5
	github.com/distribution/distribution v2.8.1+incompatible
	github.com/google/go-containerregistry v0.10.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.4.0
	gotest.tools/v3 v3.1.0 // indirect
)

replace google.golang.org/grpc/naming => github.com/xiegeo/grpc-naming v1.29.1-alpha
