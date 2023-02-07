FROM veinmind/go1.18:1.5.3-stretch as builder
WORKDIR /build
COPY . .
RUN make build

FROM dvkunion/veinmind-basic:latest as basic
FROM dvkunion/veinmind-malicious:latest as malicious
FROM dvkunion/veinmind-weakpass:latest as weakpass
# update sensitive
FROM dvkunion/veinmind-sensitive:latest as sensitive
FROM dvkunion/veinmind-history:latest as history
FROM dvkunion/veinmind-backdoor:latest as backdoor
FROM dvkunion/veinmind-log4j2:latest as log4j2
FROM dvkunion/veinmind-unsafe-mount:latest as unsafe-mount
FROM dvkunion/veinmind-iac:latest as iac
FROM dvkunion/veinmind-vuln:latest as vuln
FROM dvkunion/veinmind-escalate:latest as escalate

FROM veinmind/python3.6:1.5.3-stretch as release
WORKDIR /tool

COPY dockerfiles/sources.list /etc/apt/sources.list
RUN apt update && apt-get install -y clamav clamav-daemon && freshclam
COPY dockerfiles/clamd.conf /etc/clamav/clamd.conf

COPY --from=builder /build/veinmind-runner .
COPY --from=basic /tool/veinmind-basic .
COPY --from=weakpass /tool/veinmind-weakpass .
COPY --from=malicious /tool/veinmind-malicious_linux_amd64 .
COPY --from=sensitive /tool /tool/veinmind-sensitive
COPY --from=history /tool /tool/veinmind-history
COPY --from=backdoor /tool /tool/veinmind-backdoor
COPY --from=log4j2 /tool/veinmind-log4j2  .
COPY --from=unsafe-mount /tool/veinmind-unsafe-mount .
COPY --from=iac /tool/veinmind-iac .
COPY --from=vuln /tool/veinmind-vuln .
COPY --from=escalate /tool/veinmind-escalate .

RUN pip install -r veinmind-history/requirements.txt && pip install -r veinmind-backdoor/requirements.txt && chmod +x */scan.py
RUN echo "#!/bin/bash\n\n./veinmind-runner \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-runner
ENTRYPOINT ["/tool/entrypoint.sh"]


