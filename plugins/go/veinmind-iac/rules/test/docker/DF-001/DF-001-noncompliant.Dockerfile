# noncompliant: exposed port 22
FROM scratch
USER notroot

EXPOSE 22
HEALTHCHECK CMD curl -f http://localhost/ || exit 1