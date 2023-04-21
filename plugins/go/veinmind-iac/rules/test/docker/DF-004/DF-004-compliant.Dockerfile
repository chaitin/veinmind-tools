# compliant: `--platform` is not used to specify the mirroring platform
FROM scratch
USER notroot

EXPOSE 21
HEALTHCHECK CMD curl -f http://localhost/ || exit 1