.DEFAULT_GOAL := help

# plugin params
SHELL = /bin/bash
APP = veinmind-privilege-escalation
CMD = cmd/cli.go
ARG = ""

# build params
IMAGE_TAG = latest
CI_GOOS = linux
CI_GOARCH=$(shell uname -m)

##@ Init
.PHONY: deps
deps: 	## Install Dependencies.
ifeq ($(strip $(GITHUB_ACTIONS)),)
		go env -w GOPROXY=https://goproxy.cn,direct
endif
	go mod tidy


##@ Build
.PHONY: build
build: deps  ## Build Apps.
	go build -ldflags '-s -w' -v -trimpath -a -o ${APP} ${CMD}

.PHONY: build.platform
build.platform: deps  ## Build Apps With Platform.
	export CGO_ENABLED=1 GOOS="${CI_GOOS}" GOARCH="${CI_GOARCH}"; \
	go build -ldflags '-s -w' -v -trimpath -a -o ${APP}_${CI_GOOS}_${CI_GOARCH} ${CMD}

.PHONY: build.docker
build.docker: ## Build Apps Docker Images.
	docker build -t ${APP}:${IMAGE_TAG} .

##@ Run
.PHONY: run
run: deps ## Run Apps. 			e.g. : `make run ARG="scan iac xxxx"` .
	go run ${CMD} ${ARG}

.PHONY: run.docker
run.docker:	## Run With Parallel Container Mode. e.g. : `make run.docker ARG="scan iac xxxx"` .
	docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' -v `pwd`:/tool/data registry.veinmind.tech/veinmind/${APP} ${ARG}

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-\\.]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
