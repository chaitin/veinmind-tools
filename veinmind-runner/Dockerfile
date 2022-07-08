FROM veinmind/go1.18:1.1.0-stretch as builder
WORKDIR /build
COPY . .
RUN chmod +x script/build.sh && /bin/bash script/build.sh

# init the veinmind plugins
FROM veinmind/veinmind-malicious:latest as malicious
FROM veinmind/veinmind-weakpass:latest as weakpass
FROM veinmind/veinmind-sensitive:latest as sensitive
FROM veinmind/veinmind-history:latest as history
FROM veinmind/veinmind-backdoor:latest as backdoor

FROM veinmind/python3:1.1.0-stretch as release
WORKDIR /tool

# init the clamAV env
COPY dockerfiles/sources.list /etc/apt/sources.list
RUN apt update && apt-get install -y clamav clamav-daemon clamav-freshclam && freshclam
COPY dockerfiles/clamd.conf /etc/clamav/clamd.conf
COPY dockerfiles/clamd.sh .

# from the veinmind plugins images get tools
COPY --from=builder /build/veinmind-runner .
COPY --from=weakpass /tool/veinmind-weakpass .
COPY --from=malicious /tool/veinmind-malicious_linux_amd64 ./veinmind-malicious
COPY --from=sensitive /tool /tool/veinmind-sensitive
COPY --from=history /tool /tool/veinmind-history
COPY --from=backdoor /tool /tool/veinmind-backdoor

RUN pip install -r veinmind-sensitive/requirements.txt && pip install -r veinmind-history/requirements.txt && pip install -r veinmind-backdoor/requirements.txt && chmod +x */scan.py
RUN echo "#!/bin/bash\n\n/bin/bash clamd.sh\n\n./veinmind-runner \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-runner
ENTRYPOINT ["/tool/entrypoint.sh"]