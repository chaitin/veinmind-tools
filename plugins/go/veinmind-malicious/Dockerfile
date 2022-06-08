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
RUN echo "#!/bin/bash\n\n/bin/bash clamd.sh\n\n./veinmind-malicious_${CI_GOOS}_${CI_GOARCH} \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-malicious_${CI_GOOS}_${CI_GOARCH}
ENTRYPOINT ["/tool/entrypoint.sh"]

