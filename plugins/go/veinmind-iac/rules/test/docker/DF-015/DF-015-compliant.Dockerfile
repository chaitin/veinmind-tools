# compliant: Use Absolute Path
FROM scratch
USER notroot

ADD /root/1.tar 1.tar
HEALTHCHECK CMD curl -f http://localhost/ || exit 1