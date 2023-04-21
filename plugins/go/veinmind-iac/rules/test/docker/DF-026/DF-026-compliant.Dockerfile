# compliant: Use apk add --no-cache to clean
FROM scratch
USER notroot

RUN apk add -y --no-cache tools
HEALTHCHECK CMD curl -f http://localhost/ || exit 1