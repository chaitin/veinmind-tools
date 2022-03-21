#!/bin/bash

export GOPROXY=https://goproxy.io,direct
go mod tidy
go build -a -o veinmind-runner ./cmd/cli.go
