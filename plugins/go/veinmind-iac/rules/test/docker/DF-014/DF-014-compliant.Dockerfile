# compliant: install with auto_confirm
FROM scratch
USER notroot

RUN cat /etc/passwd
HEALTHCHECK CMD curl -f http://localhost/ || exit 1