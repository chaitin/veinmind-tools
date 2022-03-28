FROM veinmind/go1.16:1.0.2-stretch
WORKDIR /app
ADD . .
RUN export CGO_ENABLED=1 CGO_LDFLAGS_ALLOW='-Wl,.*' && \
    export GOPROXY=https://proxy.golang.com.cn,direct && \
    go build -a -tags community -o veinmind-malicious ./cmd/scan/cli.go

