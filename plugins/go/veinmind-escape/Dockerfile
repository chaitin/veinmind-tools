FROM veinmind/go1.18:1.9.42-bullseye as builder
WORKDIR /build
COPY . .
ARG TARGETARCH
ENV GOOS=linux GOARCH=${TARGETARCH}
RUN make build

FROM alpine:3.9 as compresser
WORKDIR /build
COPY --from=builder /build/veinmind-escape .
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk add upx && upx -9 veinmind-escape

FROM veinmind/base:1.9.42-bullseye as release
WORKDIR /tool
COPY --from=compresser /build/veinmind-escape .
RUN echo "#!/bin/bash\n\n./veinmind-escape \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-escape
ENTRYPOINT ["/tool/entrypoint.sh"]
