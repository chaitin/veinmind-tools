# compliant: install with auto_confirm
FROM scratch
USER notroot

WORKDIR /tool
COPY --from=release --link /build/veinmind-iac .
HEALTHCHECK CMD curl -f http://localhost/ || exit 1
RUN RUN apt-get install -y tools
ENTRYPOINT ["/tool/entrypoint.sh"]
