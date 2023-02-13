module github.com/chaitin/veinmind-tools/veinmind-runner

go 1.16

require (
	github.com/BurntSushi/toml v1.1.0
	github.com/chaitin/libveinmind v1.5.2
	github.com/chaitin/veinmind-common-go v1.2.5
	github.com/distribution/distribution v2.8.1+incompatible
	github.com/docker/docker v20.10.20+incompatible
	github.com/fatih/color v1.13.0
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/gin-gonic/gin v1.8.1
	github.com/go-git/go-git/v5 v5.4.2
	github.com/gogf/gf v1.16.9
	github.com/google/go-containerregistry v0.12.1
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/jedib0t/go-pretty/v6 v6.4.4
	github.com/matryer/is v1.4.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/rs/xid v1.4.0
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/cobra v1.6.0
	github.com/spf13/viper v1.10.0 // indirect
	go.opentelemetry.io/otel v1.7.0 // indirect
	golang.org/x/sync v0.1.0
)

replace google.golang.org/grpc/naming => github.com/xiegeo/grpc-naming v1.29.1-alpha
