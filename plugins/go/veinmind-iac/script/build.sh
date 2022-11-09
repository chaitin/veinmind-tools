#!/bin/bash

go env -w GOPROXY=https://goproxy.io,direct
go mod tidy
go build -ldflags="-w -s"  -a -o veinmind-iac ./cmd/cli.go
