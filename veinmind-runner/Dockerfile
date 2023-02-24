FROM veinmind/go1.18:1.5.3-stretch as builder
WORKDIR /build
COPY . .
ARG TARGETARCH
ENV GOOS=linux GOARCH=${TARGETARCH}
RUN make build

FROM registry.veinmind.tech/veinmind/veinmind-basic:latest as basic
FROM registry.veinmind.tech/veinmind/veinmind-malicious:latest as malicious
FROM registry.veinmind.tech/veinmind/veinmind-weakpass:latest as weakpass
# update sensitive
FROM registry.veinmind.tech/veinmind/veinmind-sensitive:latest as sensitive
FROM registry.veinmind.tech/veinmind/veinmind-history:latest as history
FROM registry.veinmind.tech/veinmind/veinmind-backdoor:latest as backdoor
FROM registry.veinmind.tech/veinmind/veinmind-log4j2:latest as log4j2
FROM registry.veinmind.tech/veinmind/veinmind-unsafe-mount:latest as unsafe-mount
FROM registry.veinmind.tech/veinmind/veinmind-iac:latest as iac
FROM registry.veinmind.tech/veinmind/veinmind-vuln:latest as vuln
FROM registry.veinmind.tech/veinmind/veinmind-escalate:latest as escalate

FROM veinmind/python3.6:1.5.3-stretch as release
WORKDIR /tool

COPY dockerfiles/sources.list /etc/apt/sources.list
RUN apt update && apt-get install -y --no-install-recommends clamav clamav-daemon && \
    rm -rf /var/lib/apt/lists/* && mkdir plugin && freshclam
COPY dockerfiles/clamd.conf /etc/clamav/clamd.conf

COPY --from=builder /build/veinmind-runner .
COPY --from=basic /tool/veinmind-basic ./plugin/
COPY --from=weakpass /tool/veinmind-weakpass ./plugin/
COPY --from=malicious /tool/veinmind-malicious ./plugin/
COPY --from=sensitive /tool/veinmind-sensitive ./plugin/
COPY --from=history /tool /tool/plugin/veinmind-history
COPY --from=backdoor /tool /tool/plugin/veinmind-backdoor
COPY --from=log4j2 /tool/veinmind-log4j2  ./plugin/
COPY --from=unsafe-mount /tool/veinmind-unsafe-mount ./plugin
COPY --from=iac /tool/veinmind-iac ./plugin
COPY --from=vuln /tool/veinmind-vuln ./plugin
COPY --from=escalate /tool/veinmind-escalate ./plugin

RUN pip install -r plugin/veinmind-history/requirements.txt && pip install -r plugin/veinmind-backdoor/requirements.txt && chmod +x plugin/*/scan.py
RUN echo "#!/bin/bash\n\n./veinmind-runner \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-runner
ENTRYPOINT ["/tool/entrypoint.sh"]


