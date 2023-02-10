FROM veinmind/go1.18:1.5.3-stretch as builder
WORKDIR /build
COPY . .
RUN make build

FROM alpine:3.9 as compresser
WORKDIR /build
COPY --from=builder /build/veinmind-malicious .
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk add upx && upx -9 veinmind-malicious

FROM veinmind/base:1.5.3-stretch
# install clamav
COPY dockerfiles/sources.list /etc/apt/sources.list
RUN apt update && apt-get install -y clamav clamav-daemon && freshclam
COPY dockerfiles/clamd.conf /etc/clamav/clamd.conf

# copy veinmind-malicious from context
ARG CI_GOOS
ENV CI_GOOS $CI_GOOS
ARG CI_GOARCH
ENV CI_GOARCH $CI_GOARCH
WORKDIR /tool
COPY --from=compresser /build/veinmind-malicious ./veinmind-malicious_${CI_GOOS}_${CI_GOARCH}
RUN echo "#!/bin/bash\n\n/bin/bash clamd\n\n./veinmind-malicious_${CI_GOOS}_${CI_GOARCH} \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-malicious_${CI_GOOS}_${CI_GOARCH}
ENTRYPOINT ["/tool/entrypoint.sh"]

