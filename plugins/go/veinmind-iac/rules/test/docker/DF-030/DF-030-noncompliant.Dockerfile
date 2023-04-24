# noncompliant: useradd without --no-log-init
FROM scratch
USER notroot

RUN useradd test
HEALTHCHECK CMD curl -f http://localhost/ || exit 1
ENTRYPOINT ["/tool/entrypoint.sh"]
