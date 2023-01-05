FROM veinmind/go1.18:1.5.1-stretch as builder
WORKDIR /build
COPY . .
RUN chmod +x script/build.sh && /bin/bash script/build.sh

FROM veinmind/base:1.5.1-stretch as release
WORKDIR /tool
COPY --from=builder /build/veinmind-sensitive .
RUN echo "#!/bin/bash\n\n./veinmind-sensitive \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-sensitive
ENTRYPOINT ["/tool/entrypoint.sh"]
