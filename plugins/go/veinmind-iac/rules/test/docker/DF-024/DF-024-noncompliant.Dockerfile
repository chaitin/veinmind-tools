# noncompliant: using update alone
FROM scratch
USER notroot

FROM veinmind/base:1.5.3-slim as compressor
WORKDIR /tool
COPY --from=release --link /build/veinmind-iac .
HEALTHCHECK CMD curl -f http://localhost/ || exit 1
RUN apt-get update
ENTRYPOINT ["/tool/entrypoint.sh"]
