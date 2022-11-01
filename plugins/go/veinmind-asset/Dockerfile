FROM veinmind/go1.18:1.3.1-stretch as builder
WORKDIR /build
COPY . .
RUN chmod +x script/build.sh && /bin/bash script/build.sh

FROM veinmind/base:1.3.1-stretch as release
WORKDIR /tool
COPY --from=builder /build/veinmind-asset .
RUN echo "#!/bin/bash\n\n./veinmind-asset \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-asset
ENTRYPOINT ["/tool/entrypoint.sh"]
