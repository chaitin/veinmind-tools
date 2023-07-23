FROM veinmind/go1.18:1.9.42-bullseye as builder
WORKDIR /build
COPY .. .
RUN make build

FROM alpine:3.9 as compresser
WORKDIR /build
COPY --from=builder /build/veinmind-trace .
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk add upx && upx -9 veinmind-trace

FROM veinmind/base:1.9.42-bullseye as release
WORKDIR /tool
COPY --from=compresser /build/veinmind-trace .
RUN echo "#!/bin/bash\n\n./veinmind-trace \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-trace
ENTRYPOINT ["/tool/entrypoint.sh"]
