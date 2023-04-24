# compliant:  Using ARG before FROM
ARG test=slim
FROM veinmind/base:1.5.3-${test} as compressor
USER notroot
WORKDIR /${test}
COPY --from=release --link /build/veinmind-iac .
RUN echo "#!/bin/bash\n\n./veinmind-iac \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-iac
HEALTHCHECK CMD curl -f http://localhost/ || exit 1
ENTRYPOINT ["/tool/entrypoint.sh"]
