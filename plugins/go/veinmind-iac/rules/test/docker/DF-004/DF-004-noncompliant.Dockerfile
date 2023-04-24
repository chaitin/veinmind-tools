# noncompliant: `--platform` is used to specify the mirroring platform
FROM --platform=arm64 scratch
USER notroot

EXPOSE 21
HEALTHCHECK CMD curl -f http://localhost/ || exit 1
