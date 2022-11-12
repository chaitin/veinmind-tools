#!/bin/bash

go env -w GOPROXY=https://goproxy.io,direct
go mod tidy
go build -ldflags '-s -w' -trimpath -a -o veinmind-log4j2 ./cmd/cli.go
