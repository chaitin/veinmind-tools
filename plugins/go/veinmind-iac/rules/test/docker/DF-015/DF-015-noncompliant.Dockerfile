# noncompliant: Use add with ../ at src dir
FROM scratch
USER notroot

ADD ../1.tar 1.tar
HEALTHCHECK CMD curl -f http://localhost/ || exit 1