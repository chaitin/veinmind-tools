FROM veinmind/go1.16:1.0.2-stretch as builder
WORKDIR /build
COPY . .
RUN chmod +x script/build.sh && /bin/bash script/build.sh

FROM veinmind/base:1.0.2-stretch as release
WORKDIR /tool
COPY --from=builder /build/veinmind-runner .
RUN echo "#!/bin/bash\n\n./veinmind-runner \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-runner
ENTRYPOINT ["/tool/entrypoint.sh"]


