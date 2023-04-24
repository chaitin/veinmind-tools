# compliant: After using `RUN yum install` there is a cleanup
FROM scratch
USER notroot

WORKDIR /tool
RUN yum install -y tools && yum clean all
HEALTHCHECK CMD curl -f http://localhost/ || exit 1
ENTRYPOINT ["/tool/entrypoint.sh"]