# compliant: Do not exposed port 22
FROM scratch
USER notroot

EXPOSE 21
HEALTHCHECK CMD curl -f http://localhost/ || exit 1