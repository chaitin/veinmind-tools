# noncompliant: install missing auto_confirm
FROM scratch
USER notroot

WORKDIR /tool
COPY --from=release --link /build/veinmind-iac .
HEALTHCHECK CMD curl -f http://localhost/ || exit 1
RUN apt-get install tools
ENTRYPOINT ["/tool/entrypoint.sh"]