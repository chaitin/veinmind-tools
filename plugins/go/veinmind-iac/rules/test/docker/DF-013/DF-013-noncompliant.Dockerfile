# noncompliant: use curl and wget
FROM scratch
USER notroot

RUN wget http://localhost:8888 && curl http://localhost:8889
HEALTHCHECK CMD curl -f http://localhost/ || exit 1