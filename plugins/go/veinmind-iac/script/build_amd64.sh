go mod tidy
mkdir -p ./artifacts/${CI_GOOS}-${CI_GOARCH}
export GOOS="$CI_GOOS" GOARCH="$CI_GOARCH"
go build -ldflags '-s -w' -trimpath -a -o ./artifacts/${CI_GOOS}-${CI_GOARCH}/veinmind-iac_${CI_GOOS}_${CI_GOARCH} ./cmd/cli.go
