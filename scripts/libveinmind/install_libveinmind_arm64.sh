#!/bin/bash

curl -O http://download.veinmind.tech/libveinmind/dists%2Fmain/libveinmind-dev_1.0.1-1_arm64.deb
sudo dpkg-deb -X libveinmind-dev_1.0.1-1_arm64.deb /