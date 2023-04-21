# noncompliant: missing clean after using apk add
FROM scratch
USER notroot

RUN apk add -y tools
HEALTHCHECK CMD curl -f http://localhost/ || exit 1