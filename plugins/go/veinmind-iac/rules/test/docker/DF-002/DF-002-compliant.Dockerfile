# compliant: `EXPOSE` does not specifie an out-of-range port
FROM scratch
USER notroot

EXPOSE 21
HEALTHCHECK CMD curl -f http://localhost/ || exit 1