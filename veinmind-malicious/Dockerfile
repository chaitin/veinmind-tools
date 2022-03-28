FROM veinmind/go1.16:1.0.2-stretch AS build
WORKDIR /app
ADD . .
RUN export CGO_ENABLED=1 CGO_LDFLAGS_ALLOW='-Wl,.*' && \
    export GOPROXY=https://proxy.golang.com.cn,direct && \
    go build -a -tags community -o veinmind-malicious ./cmd/scan/cli.go
FROM veinmind/base:1.0.2-stretch
WORKDIR /tool
COPY --from=build /app/veinmind-malicious /tool/veinmind-malicious
RUN echo "#!/bin/bash\n\n./veinmind-malicious \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-malicious
ENTRYPOINT ["/tool/entrypoint.sh"]

