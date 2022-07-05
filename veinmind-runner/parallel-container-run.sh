#!/bin/bash

docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' -v `pwd`:/tool/resource -v /var/run/docker.sock:/var/run/docker.sock veinmind/veinmind-runner $*
