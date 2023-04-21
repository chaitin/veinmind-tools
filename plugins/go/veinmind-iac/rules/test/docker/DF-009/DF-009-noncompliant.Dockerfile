# noncompliant: The chown tag exists for the copy command
FROM scratch
USER notroot

COPY --chown=1 --link aa /
HEALTHCHECK CMD curl -f http://localhost/ || exit 1