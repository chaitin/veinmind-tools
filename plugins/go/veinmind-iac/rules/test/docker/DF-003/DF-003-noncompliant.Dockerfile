# noncompliant:  Using latest images
FROM veinmind/base:1.5.3-slim:latest as compressor
USER notroot
WORKDIR /tool
COPY --from=release --link /build/veinmind-iac .
RUN echo "#!/bin/bash\n\n./veinmind-iac \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-iac
HEALTHCHECK CMD curl -f http://localhost/ || exit 1
ENTRYPOINT ["/tool/entrypoint.sh"]