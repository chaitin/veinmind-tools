# noncompliant: After using `RUN zypper in/remove/source-install.... ' without cleaning up after
FROM scratch
USER notroot

WORKDIR /tool
RUN zypper -y install tools
ENTRYPOINT ["/tool/entrypoint.sh"]
HEALTHCHECK CMD curl -f http://localhost/ || exit 1