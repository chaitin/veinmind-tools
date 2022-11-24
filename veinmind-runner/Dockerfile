FROM veinmind/go1.18:1.3.1-stretch as builder
WORKDIR /build
COPY . .
RUN chmod +x script/build.sh && /bin/bash script/build.sh

# FROM veinmind/veinmind-malicious:latest as malicious
FROM veinmind/veinmind-weakpass:latest as weakpass
FROM veinmind/veinmind-sensitive:latest as sensitive
FROM veinmind/veinmind-history:latest as history
FROM veinmind/veinmind-backdoor:latest as backdoor
FROM veinmind/veinmind-log4j2:latest as log4j2
FROM veinmind/veinmind-unsafe-mount:latest as unsafe-mount
FROM veinmind/veinmind-iac:latest as iac

FROM veinmind/python3.6:1.3.1-stretch as release
WORKDIR /tool

COPY dockerfiles/sources.list /etc/apt/sources.list
RUN apt update && apt-get install -y clamav clamav-daemon && freshclam
COPY dockerfiles/clamd.conf /etc/clamav/clamd.conf

COPY --from=builder /build/veinmind-runner .
COPY --from=weakpass /tool/veinmind-weakpass .
# COPY --from=malicious /tool/veinmind-malicious_linux_amd64 .
COPY --from=sensitive /tool /tool/veinmind-sensitive
COPY --from=history /tool /tool/veinmind-history
COPY --from=backdoor /tool /tool/veinmind-backdoor
COPY --from=log4j2 /tool/veinmind-log4j2  .
COPY --from=unsafe-mount /tool/veinmind-unsafe-mount .
COPY --from=iac /tool/veinmind-iac .

RUN pip install -r veinmind-sensitive/requirements.txt && pip install -r veinmind-history/requirements.txt && pip install -r veinmind-backdoor/requirements.txt && chmod +x */scan.py
RUN echo "#!/bin/bash\n\n./veinmind-runner \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-runner
ENTRYPOINT ["/tool/entrypoint.sh"]


