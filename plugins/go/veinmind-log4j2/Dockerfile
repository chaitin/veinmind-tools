FROM veinmind/go1.18:1.3.1-stretch as builder
WORKDIR /build
COPY . .
RUN chmod +x script/build.sh && /bin/bash script/build.sh

FROM veinmind/base:1.3.1-stretch as release
WORKDIR /tool
COPY --from=builder /build/veinmind-log4j2 .
RUN echo "#!/bin/bash\n\n./veinmind-log4j2 \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-log4j2
ENTRYPOINT ["/tool/entrypoint.sh"]
