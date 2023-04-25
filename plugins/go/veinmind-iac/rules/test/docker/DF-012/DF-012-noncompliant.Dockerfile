# noncompliant: Sensitive information in env
FROM scratch
USER notroot

WORKDIR /tool
ENV passwd qqq
HEALTHCHECK CMD curl -f http://localhost/ || exit 1