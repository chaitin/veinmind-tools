FROM veinmind/go1.18:1.5.3-stretch as builder
WORKDIR /build
COPY . .
RUN make build

FROM alpine:3.9 as compresser
WORKDIR /build
COPY --from=builder /build/veinmind-vuln .
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk add upx && upx -9 veinmind-vuln

FROM veinmind/base:1.5.3-slim as release
WORKDIR /tool
COPY --from=compresser /build/veinmind-vuln .
RUN echo "#!/bin/bash\n\n./veinmind-vuln \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-vuln
ENTRYPOINT ["/tool/entrypoint.sh"]
