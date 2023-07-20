FROM veinmind/go1.18:1.9.42-bullseye as builder
WORKDIR /build
COPY . .
ARG TARGETARCH
ENV GOOS=linux GOARCH=${TARGETARCH}
RUN make build

FROM alpine:3.9 as compresser
WORKDIR /build
COPY --from=builder /build/veinmind-malicious .
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk add upx && upx -9 veinmind-malicious

FROM veinmind/base:1.9.42-bullseye as release
WORKDIR /tool
COPY dockerfiles/sources.list /etc/apt/sources.list
RUN apt update && apt-get install -y --no-install-recommends --allow-remove-essential libtinfo5 clamav clamav-daemon && rm -rf /var/lib/apt/lists/* && freshclam
COPY dockerfiles/clamd.conf /etc/clamav/clamd.conf
COPY --from=compresser /build/veinmind-malicious ./veinmind-malicious
RUN echo "#!/bin/bash\n\n/bin/bash -c /usr/sbin/clamd -F --config-file=/etc/clamav/clamd.conf\n\n./veinmind-malicious \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-malicious
ENTRYPOINT ["/tool/entrypoint.sh"]



