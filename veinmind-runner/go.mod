module github.com/chaitin/veinmind-tools/veinmind-runner

go 1.16

replace github.com/chaitin/veinmind-tools/veinmind-common/go v1.0.0 => ../veinmind-common/go

require (
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/chaitin/libveinmind v1.0.4
	github.com/chaitin/veinmind-tools/veinmind-common/go v1.0.0
	github.com/containerd/containerd v1.6.1 // indirect
	github.com/docker/cli v20.10.12+incompatible
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.13+incompatible
	github.com/fvbommel/sortorder v1.0.2 // indirect
	github.com/google/go-containerregistry v0.8.0
	github.com/moby/sys/mount v0.3.1 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/spf13/cobra v1.3.0
	github.com/theupdateframework/notary v0.7.0 // indirect
	gotest.tools/v3 v3.1.0 // indirect
)
