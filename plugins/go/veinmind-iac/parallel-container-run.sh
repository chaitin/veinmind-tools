#!/bin/bash

docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' -v `pwd`:/tool/data veinmind/veinmind-iac $*