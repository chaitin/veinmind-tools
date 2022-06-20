FROM veinmind/go1.18:1.1.0-stretch as builder
WORKDIR /build
COPY . .
RUN chmod +x script/build.sh && /bin/bash script/build.sh

FROM veinmind/base:1.1.0-stretch as release
WORKDIR /tool
COPY --from=builder /build/veinmind-weakpass .
RUN echo "#!/bin/bash\n\n./veinmind-weakpass \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-weakpass
ENTRYPOINT ["/tool/entrypoint.sh"]

