#!/bin/bash

export GOPROXY=https://goproxy.io,direct
go mod tidy
go build -ldflags '-s -w' -trimpath -a -o veinmind-runner ./cmd/
