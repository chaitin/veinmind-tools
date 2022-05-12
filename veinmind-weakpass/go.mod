module github.com/chaitin/veinmind-tools/veinmind-weakpass

go 1.17

replace github.com/chaitin/veinmind-tools/veinmind-common/go v1.0.0 => ../veinmind-common/go

require (
	github.com/Jeffail/tunny v0.1.4
	github.com/chaitin/libveinmind v1.0.7
	github.com/chaitin/veinmind-tools/veinmind-common/go v1.0.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.3.0
	golang.org/x/crypto v0.0.0-20220112180741-5e0467b6c7ce
)

require (
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20211205182925-97ca703d548d // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)
