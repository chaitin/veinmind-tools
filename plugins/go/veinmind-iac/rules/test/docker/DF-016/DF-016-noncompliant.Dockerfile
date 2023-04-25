# noncompliant: Multiple CMD instructions listed
FROM scratch
USER notroot

CMD curl -f http://localhost/ || exit 2
CMD curl -f http://localhost/ || exit 1
HEALTHCHECK CMD curl -f http://localhost/ || exit 1