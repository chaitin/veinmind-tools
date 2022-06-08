#!/bin/bash

go env -w GOPROXY=https://goproxy.io,direct
go mod tidy
go build -a -o veinmind-asset ./cmd/cli.go
