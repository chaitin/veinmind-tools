# noncompliant:  `COPY --from` is from the current image
FROM scratch as release
USER notroot

COPY --from=release --link /build/veinmind-iac .
RUN echo "#!/bin/bash\n\n./veinmind-iac \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-iac
HEALTHCHECK CMD curl -f http://localhost/ || exit 1
ENTRYPOINT ["/tool/entrypoint.sh"]

