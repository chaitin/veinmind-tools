#!/bin/bash

# use: bash scripts/libveinmind/upgrade_libveinmind_version.sh YOUR_VERSION
# example: bash scripts/libveinmind/upgrade_libveinmind_version.sh 1.3.1
# can not use script with version not like x.x.x

# if you want to upgrade libveinmind, you should upgrade three part
version=$1
# check darwin or linux
shopt -s expand_aliases
if [[ $(uname) == "Darwin" ]]; then
  echo "darwin"
  alias sed='sed -i ""'
elif [[ $(uname) == "Linux" ]]; then
  echo "linux"
  alias sed='sed -i'
fi

# part.1: go plugins(extends veinmind-runner) go.mod
sed "s/github\.com\/chaitin\/libveinmind v[0-9]\.[0-9]\.[0-9]$/github.com\/chaitin\/libveinmind v${version}/g" $(grep "github.com/chaitin/libveinmind v[0-9]\.[0-9]\.[0-9]$" -rl ./plugins ./example)

# part.2: python plugins requirements.txt
sed "s/veinmind==[0-9]\.[0-9]\.[0-9]/veinmind==${version}/g" $(grep "veinmind==" -rl ./plugins ./example)
#
## part.3: image tag in Dockerfile/github-workflow
sed "s/veinmind\/python3[0-9\.]*:[0-9]\.[0-9]\.[0-9]/veinmind\/python3.6:${version}/g" $(grep "veinmind\/python3[0-9\.]*:[0-9]\.[0-9]\.[0-9]" -rl ./plugins ./example ./veinmind-runner ./.github)
sed "s/veinmind\/go1.*:[0-9]\.[0-9]\.[0-9]/veinmind\/go1.18:${version}/g" $(grep "veinmind\/go1.*:[0-9]\.[0-9]\.[0-9]" -rl ./plugins ./example ./veinmind-runner ./.github)
sed "s/veinmind\/base:[0-9]\.[0-9]\.[0-9]/veinmind\/base:${version}/g" $(grep "veinmind\/base:[0-9]\.[0-9]\.[0-9]" -rl ./plugins ./example ./veinmind-runner ./.github)
