FROM veinmind/go1.18:1.5.3-stretch as builder
WORKDIR /build
COPY . .
RUN make build

FROM alpine:3.9 as compresser
WORKDIR /build
COPY --from=builder /build/veinmind-weakpass .
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk add upx && upx -9 veinmind-weakpass

FROM veinmind/base:1.5.3-slim as release
WORKDIR /tool
COPY --from=compresser /build/veinmind-weakpass .
RUN echo "#!/bin/bash\n\n./veinmind-weakpass \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-weakpass
ENTRYPOINT ["/tool/entrypoint.sh"]

