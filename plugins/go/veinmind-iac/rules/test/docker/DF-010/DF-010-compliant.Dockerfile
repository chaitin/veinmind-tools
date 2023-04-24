# compliant: use copy instead add
FROM scratch
USER notroot

COPY --link 1.py 2.py
HEALTHCHECK CMD curl -f http://localhost/ || exit 1