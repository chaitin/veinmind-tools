# noncompliant:  Multiple ENTRYPOINT instructions listed
FROM scratch
USER notroot

RUN echo "#!/bin/bash\n\n./veinmind-iac \$*" > /tool/entrypoint.sh && chmod +x /tool/entrypoint.sh && chmod +x /tool/veinmind-iac
HEALTHCHECK CMD curl -f http://localhost/ || exit 1
ENTRYPOINT ["/tool/entrypoint.sh"]
ENTRYPOINT ["/tool/entrypoint1.sh"]
