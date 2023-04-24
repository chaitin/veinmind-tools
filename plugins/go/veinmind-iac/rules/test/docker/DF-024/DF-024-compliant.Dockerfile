# compliant: using without sudo
FROM scratch
USER notroot

FROM veinmind/base:1.5.3-slim as compressor
WORKDIR /tool
COPY --from=release --link /build/veinmind-iac .
RUN apt-get update && apt-get install -y tools && apt-get clean
HEALTHCHECK CMD curl -f http://localhost/ || exit 1
ENTRYPOINT ["/tool/entrypoint.sh"]