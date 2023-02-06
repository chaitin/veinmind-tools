.PHONY: deps build upx

SHELL = /bin/bash
APP = veinmind-basic
CMD = cmd/basic/cli.go

deps:
	go env -w GOPROXY=https://goproxy.io,direct && go mod tidy

build: deps
	go build -ldflags '-s -w' -v -trimpath -a -o ${APP} ${CMD}