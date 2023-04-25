# compliant: use RUN curl/wget instead ADD
FROM scratch
USER notroot

RUN wget https://www.download.com/1.tar 1.tar
HEALTHCHECK CMD curl -f http://localhost/ || exit 1