FROM veinmind/go1.18:1.9.42-bullseye as builder
WORKDIR /build
COPY . .
ARG TARGETARCH
ENV GOOS=linux GOARCH=${TARGETARCH}
RUN make build

FROM veinmind/veinmind-basic:latest as basic
FROM veinmind/veinmind-malicious:latest as malicious
FROM veinmind/veinmind-weakpass:latest as weakpass
# update sensitive
FROM veinmind/veinmind-sensitive:latest as sensitive
FROM veinmind/veinmind-history:latest as history
FROM veinmind/veinmind-backdoor:latest as backdoor
FROM veinmind/veinmind-log4j2:latest as log4j2
FROM veinmind/veinmind-unsafe-mount:latest as unsafe-mount
FROM veinmind/veinmind-iac:latest as iac
FROM veinmind/veinmind-vuln:latest as vuln
FROM veinmind/veinmind-escape:latest as escape
FROM veinmind/veinmind-privilege-escalation:latest as privilege-escalation
FROM veinmind/veinmind-minio:latest as minio
FROM veinmind/veinmind-trace:latest as trace

FROM veinmind/python3.6:1.9.42-bullseye as release
WORKDIR /tool

COPY dockerfiles/sources.list /etc/apt/sources.list
RUN apt update && apt-get install -y --no-install-recommends --allow-remove-essential libtinfo5 clamav clamav-daemon && \
    rm -rf /var/lib/apt/lists/* && mkdir plugin && freshclam
COPY dockerfiles/clamd.conf /etc/clamav/clamd.conf

COPY --from=builder /build/veinmind-runner .
COPY --from=basic /tool/veinmind-basic ./plugin/
COPY --from=weakpass /tool/veinmind-weakpass ./plugin/
COPY --from=malicious /tool/veinmind-malicious ./plugin/
COPY --from=sensitive /tool/veinmind-sensitive ./plugin/
COPY --from=history /tool /tool/plugin/veinmind-history
COPY --from=backdoor /tool /tool/veinmind-backdoor
COPY --from=log4j2 /tool/veinmind-log4j2  ./plugin/
COPY --from=unsafe-mount /tool/veinmind-unsafe-mount ./plugin
COPY --from=iac /tool/veinmind-iac ./plugin
COPY --from=vuln /tool/veinmind-vuln ./plugin
COPY --from=escape /tool/veinmind-escape ./plugin
COPY --from=privilege-escalation /tool/veinmind-privilege-escalation ./plugin
COPY --from=minio /tool/veinmind-minio ./plugin
COPY --from=trace /tool/veinmind-trace ./plugin

RUN pip install -r plugin/veinmind-history/requirements.txt && chmod +x plugin/*/scan.py
RUN echo "#!/bin/bash\n\n./veinmind-runner \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-runner
ENTRYPOINT ["/tool/entrypoint.sh"]


