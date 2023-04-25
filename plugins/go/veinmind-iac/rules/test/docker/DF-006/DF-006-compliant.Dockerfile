# compliant: `WORKDIR` uses a absolute path
FROM scratch
USER notroot


WORKDIR /tool
HEALTHCHECK CMD curl -f http://localhost/ || exit 1