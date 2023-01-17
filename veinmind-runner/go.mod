module github.com/chaitin/veinmind-tools/veinmind-runner

go 1.16

require (
	github.com/BurntSushi/toml v1.1.0
	github.com/chaitin/libveinmind v1.5.2
	github.com/chaitin/veinmind-common-go v1.2.1
	github.com/distribution/distribution v2.8.1+incompatible
	github.com/docker/docker v20.10.17+incompatible
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/gin-gonic/gin v1.8.1
	github.com/go-git/go-git/v5 v5.4.2
	github.com/gogf/gf v1.16.9
	github.com/google/go-containerregistry v0.10.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/rs/xid v1.4.0
	github.com/sirupsen/logrus v1.8.1
	go.opentelemetry.io/otel v1.7.0 // indirect
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4
	golang.org/x/text v0.3.8-0.20211105212822-18b340fc7af2 // indirect
)

replace google.golang.org/grpc/naming => github.com/xiegeo/grpc-naming v1.29.1-alpha
