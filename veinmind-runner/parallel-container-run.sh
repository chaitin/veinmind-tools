#!/bin/bash

docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' -v /run/docker/plugins:/run/docker/plugins -v /var/run/docker.sock:/var/run/docker.sock veinmind/veinmind-runner $*
