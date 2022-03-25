#!/bin/bash

docker run --rm -it --network veinmind --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' -v `pwd`:/tool/data -e CLAMD_HOST=clamav -e CLAMD_PORT=3310 veinmind/veinmind-malicious $*
