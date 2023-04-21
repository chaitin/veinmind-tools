# noncompliant: useradd without --no-log-init
FROM scratch
USER notroot

RUN sudo cat /etc/passwd
HEALTHCHECK CMD curl -f http://localhost/ || exit 1