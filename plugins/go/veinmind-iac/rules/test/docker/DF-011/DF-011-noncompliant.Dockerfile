# noncompliant: use add to get https resource
FROM scratch
USER notroot

ADD https://www.download.com/1.tar 1.tar
HEALTHCHECK CMD curl -f http://localhost/ || exit 1
