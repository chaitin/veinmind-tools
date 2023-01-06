FROM veinmind/go1.18:1.5.1-stretch as builder
WORKDIR /build
COPY .. .
RUN chmod +x script/build.sh && /bin/bash script/build.sh

FROM veinmind/base:1.5.1-slim as release
WORKDIR /tool
COPY --from=builder /build/veinmind-unsafe-mount .
RUN echo "#!/bin/bash\n\n./veinmind-unsafe-mount \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-unsafe-mount
ENTRYPOINT ["/tool/entrypoint.sh"]
