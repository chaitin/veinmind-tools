sudo apt-get update && sudo apt-get install -y gcc-aarch64-linux-gnu g++-aarch64-linux-gnu
sudo apt-get install libssl-dev
export CC="aarch64-linux-gnu-gcc" CXX="aarch64-linux-gnu-g++"
export PKG_CONFIG_PATH=/usr/lib/aarch64-linux-gnu/pkgconfig
go mod tidy
mkdir -p ./artifacts/${CI_GOOS}-${CI_GOARCH}
export CGO_ENABLED=1 GOOS="$CI_GOOS" GOARCH="$CI_GOARCH"
go build -a -tags community -o ./artifacts/${CI_GOOS}-${CI_GOARCH}/veinmind-weakpass_${CI_GOOS}_${CI_GOARCH} ./cmd/scan/
