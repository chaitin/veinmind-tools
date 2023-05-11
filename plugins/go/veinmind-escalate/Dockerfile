FROM veinmind/go1.18:1.9.15-bullseye as builder
WORKDIR /build
COPY . .
ARG TARGETARCH
ENV GOOS=linux GOARCH=${TARGETARCH}
RUN make build

FROM alpine:3.9 as compresser
WORKDIR /build
COPY --from=builder /build/veinmind-escalate .
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk add upx && upx -9 veinmind-escalate

FROM veinmind/base:1.9.15-bullseye as release
WORKDIR /tool
COPY --from=compresser /build/veinmind-escalate .
RUN echo "#!/bin/bash\n\n./veinmind-escalate \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-escalate
ENTRYPOINT ["/tool/entrypoint.sh"]
