FROM veinmind/base:1.1.0-stretch
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
COPY dockerfiles/clamd.sh .
ADD veinmind-malicious_${CI_GOOS}_${CI_GOARCH} .
RUN echo "#!/bin/bash\n\n./veinmind-malicious_${CI_GOOS}_${CI_GOARCH} --clamav-manually-start=false --clamav-addr=127.0.0.1 --clamav-port=3310\$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-malicious_${CI_GOOS}_${CI_GOARCH}
ENTRYPOINT ["/tool/entrypoint.sh"]

