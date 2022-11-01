module github.com/chaitin/veinmind-tools/veinmind-runner

go 1.16

require (
	github.com/BurntSushi/toml v0.4.1
	github.com/chaitin/libveinmind v1.1.1
	github.com/chaitin/veinmind-common-go v1.1.9
	github.com/distribution/distribution v2.8.1+incompatible
	github.com/docker/docker v20.10.17+incompatible
	github.com/gin-gonic/gin v1.8.1
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.5.0
	golang.org/x/sync v0.0.0-20220513210516-0976fa681c29
	gotest.tools/v3 v3.1.0 // indirect
)

replace google.golang.org/grpc/naming => github.com/xiegeo/grpc-naming v1.29.1-alpha
