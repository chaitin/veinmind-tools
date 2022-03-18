go mod tidy
mkdir -p ./artifacts/${CI_GOOS}-${CI_GOARCH}
export CGO_ENABLED=1 GOOS="$CI_GOOS" GOARCH="$CI_GOARCH" TAGS="$TAGS"
go build -a -tags ${TAGS} -o ./artifacts/${CI_GOOS}-${CI_GOARCH}/veinmind-weakpass_${CI_GOOS}_${CI_GOARCH} ./cmd/cli.go
