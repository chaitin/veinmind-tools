#!/bin/bash

echo 'deb [trusted=yes] https://download.veinmind.tech/libveinmind/apt/ ./' | sudo tee /etc/apt/sources.list.d/libveinmind.list
sudo apt-get update
sudo apt-get install libveinmind-dev