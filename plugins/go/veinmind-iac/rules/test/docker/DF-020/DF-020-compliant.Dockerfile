# compliant:  COPY with more than two arguments end with slash
FROM scratch
USER notroot

COPY --link --from=release /build/veinmind-iac .
RUN echo "#!/bin/bash\n\n./veinmind-iac \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-iac
HEALTHCHECK CMD curl -f http://localhost/ || exit 1
ENTRYPOINT ["/tool/entrypoint.sh"]