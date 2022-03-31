module github.com/chaitin/veinmind-tools/veinmind-weakpass

go 1.16

replace github.com/chaitin/veinmind-tools/veinmind-common/go v1.0.0 => ../veinmind-common/go

require (
	github.com/Jeffail/tunny v0.1.4
	github.com/chaitin/libveinmind v1.0.4
	github.com/chaitin/veinmind-tools/veinmind-common/go v1.0.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.3.0
	golang.org/x/crypto v0.0.0-20220112180741-5e0467b6c7ce
)
