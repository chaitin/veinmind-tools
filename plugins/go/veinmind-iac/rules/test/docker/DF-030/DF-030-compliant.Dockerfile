# compliant: useradd with --no-log-init
FROM scratch
USER notroot

RUN useradd test --no-log-init
HEALTHCHECK CMD curl -f http://localhost/ || exit 1
ENTRYPOINT ["/tool/entrypoint.sh"]