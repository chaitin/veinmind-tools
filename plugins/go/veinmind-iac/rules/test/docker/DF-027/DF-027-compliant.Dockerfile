# compliant: Cleanup after using `RUN dnf install`
FROM scratch
USER notroot

HEALTHCHECK CMD curl -f http://localhost/ || exit 1
RUN dnf install tools && dnf clean all