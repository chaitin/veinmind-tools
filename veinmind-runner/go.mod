module github.com/chaitin/veinmind-tools/veinmind-runner

go 1.16

require (
	github.com/BurntSushi/toml v0.4.1
	github.com/chaitin/libveinmind v1.3.2
	github.com/chaitin/veinmind-common-go v1.2.1
	github.com/containerd/containerd v1.6.9 // indirect
	github.com/distribution/distribution v2.8.1+incompatible
	github.com/docker/docker v20.10.17+incompatible
	github.com/gin-gonic/gin v1.8.1
	github.com/go-git/go-git/v5 v5.4.2
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.5.0
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4
	gotest.tools/v3 v3.1.0 // indirect
)

replace google.golang.org/grpc/naming => github.com/xiegeo/grpc-naming v1.29.1-alpha
