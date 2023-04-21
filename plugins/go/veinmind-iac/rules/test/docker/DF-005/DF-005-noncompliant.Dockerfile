# noncompliant: user root user
FROM scratch

ENTRYPOINT ["/tool/entrypoint.sh"]
HEALTHCHECK CMD curl -f http://localhost/ || exit 1