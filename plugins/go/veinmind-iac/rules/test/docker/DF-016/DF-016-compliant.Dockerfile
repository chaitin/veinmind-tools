# compliant: only one cmd
FROM scratch
USER notroot

CMD curl -f http://localhost/ || exit 1
HEALTHCHECK CMD curl -f http://localhost/ || exit 1
