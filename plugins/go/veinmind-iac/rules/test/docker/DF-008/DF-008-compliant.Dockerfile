# compliant: only one HEALTHCHECK
FROM scratch
USER notroot

EXPOSE 21
HEALTHCHECK CMD curl https://localhost:8888