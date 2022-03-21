module github.com/chaitin/veinmind-tools/veinmind-runner

go 1.16

replace (
	github.com/chaitin/veinmind-tools/veinmind-common/go v1.0.0 => ../veinmind-common/go
)

require (
	github.com/chaitin/libveinmind v1.0.4
	github.com/sirupsen/logrus v1.8.1 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	github.com/chaitin/veinmind-tools/veinmind-common/go v1.0.0
)
