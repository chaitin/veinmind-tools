# noncompliant: `EXPOSE` specifies an out-of-range port
FROM scratch
USER notroot

EXPOSE 77777
HEALTHCHECK CMD curl -f http://localhost/ || exit 1