FROM veinmind/go1.18:1.9.42-bullseye as builder
WORKDIR /build
COPY .. .
RUN make build

FROM alpine:3.9 as compresser
WORKDIR /build
COPY --from=builder /build/veinmind-minio .
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk add upx && upx -9 veinmind-minio

FROM veinmind/base:1.9.42-bullseye as release
WORKDIR /tool
COPY --from=compresser /build/veinmind-minio .
RUN echo "#!/bin/bash\n\n./veinmind-minio \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-minio
ENTRYPOINT ["/tool/entrypoint.sh"]
