#!/bin/bash

export GOPROXY=https://goproxy.io,direct
export CGO_ENABLED=1
go mod tidy
go build -ldflags '-s -w' -trimpath -tags dynamic -a -o veinmind-weakpass ./cmd/cli.go
