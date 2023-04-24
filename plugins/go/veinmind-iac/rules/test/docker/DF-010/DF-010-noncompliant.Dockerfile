# compliant: use ADD
FROM scratch
USER notroot

ADD 1.py 2.py
HEALTHCHECK CMD curl -f http://localhost/ || exit 1