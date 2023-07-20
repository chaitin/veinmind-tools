FROM veinmind/go1.18:1.9.42-bullseye as builder
WORKDIR /build
COPY .. .
RUN make build

FROM alpine:3.9 as compresser
WORKDIR /build
COPY --from=builder /build/veinmind-privilege-escalation .
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk add upx && upx -9 veinmind-privilege-escalation

FROM veinmind/base:1.9.42-bullseye as release
WORKDIR /tool
COPY --from=compresser /build/veinmind-privilege-escalation .
RUN echo "#!/bin/bash\n\n./veinmind-privilege-escalation \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-privilege-escalation
ENTRYPOINT ["/tool/entrypoint.sh"]
