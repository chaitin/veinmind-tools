# compliant: use notroot user
FROM scratch
USER notroot

ENTRYPOINT ["/tool/entrypoint.sh"]
HEALTHCHECK CMD curl -f http://localhost/ || exit 1