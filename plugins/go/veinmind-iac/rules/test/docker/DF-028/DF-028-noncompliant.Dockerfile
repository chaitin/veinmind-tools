# noncompliant: No cleanup after using `RUN yum install`
FROM scratch
USER notroot

WORKDIR /tool
RUN yum install -y tools
HEALTHCHECK CMD curl -f http://localhost/ || exit 1
ENTRYPOINT ["/tool/entrypoint.sh"]
