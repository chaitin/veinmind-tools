# compliant: No sensitive information in env
FROM scratch
USER notroot

WORKDIR /tool
HEALTHCHECK CMD curl -f http://localhost/ || exit 1