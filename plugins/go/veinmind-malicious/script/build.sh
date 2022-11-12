#!/bin/bash

go env -w GOPROXY=https://goproxy.io,direct
go mod tidy
export CGO_ENABLED=1 CGO_LDFLAGS_ALLOW='-Wl,.*'
go build -ldflags '-s -w' -trimpath -a -o veinmind-malicious ./cmd/scan/
