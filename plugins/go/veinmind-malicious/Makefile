.DEFAULT_GOAL := help

SHELL = /bin/bash
APP = veinmind-malicious
CMD = cmd/scan/*
ARG = ""

# build params
IMAGE_TAG = latest
CI_GOOS = linux
CI_GOARCH=$(shell uname -m)
TAGS = community

##@ Init
.PHONY: deps
deps: 	## Install Dependencies.
ifeq ($(strip $(GITHUB_ACTIONS)),)
		go env -w GOPROXY=https://goproxy.cn,direct
endif
	go mod tidy

clamav: ## Install Clamav.
	apt-get update && apt-get install -y clamav clamav-daemon && freshclam

arm: ## Install arm static build .
	apt-get update && apt-get install -y gcc-aarch64-linux-gnu g++-aarch64-linux-gnu libssl-dev

##@ Build
.PHONY: build
build: deps   ## 	Build Apps.
	go build -ldflags '-s -w' -v -trimpath -a -tags ${TAGS}  -o ${APP} ${CMD}

.PHONY: build.platform
build.platform: deps ## Build Apps With Platform.
	export CGO_ENABLED=1 CGO_LDFLAGS_ALLOW='-Wl,.*' ; GOOS="${CI_GOOS}" GOARCH="${CI_GOARCH}"; \
	go build -ldflags '-s -w' -v -trimpath -a -tags ${TAGS} -o ${APP}_${CI_GOOS}_${CI_GOARCH} ${CMD}

.PHONY: build.arm
build.arm: deps arm
	export CC="aarch64-linux-gnu-gcc" CXX="aarch64-linux-gnu-g++"; \
	export PKG_CONFIG_PATH=/usr/lib/aarch64-linux-gnu/pkgconfig; \
	export CGO_ENABLED=1 CGO_LDFLAGS_ALLOW='-Wl,.*' ; GOOS="${CI_GOOS}" GOARCH="${CI_GOARCH}"; \
	go build -ldflags '-s -w' -v -trimpath -a -tags ${TAGS} -o ${APP}_${CI_GOOS}_${CI_GOARCH} ${CMD}

.PHONY: build.docker
build.docker: ## Build Apps Docker Images.
	docker build -t ${APP}:${IMAGE_TAG} .

##@ Run
.PHONY: run
run: deps run.clamav ## Run Apps. 			e.g. : `make run ARG="scan image"` .
	go run ${CMD} ${ARG}

.PHONY: run.clamav
run.clamav: deps ## Run Apps. 			e.g. : `make run ARG="scan image"` .
	cp dockerfiles/clamd.conf /etc/clamav/clamd.conf
	clamd
	go run ${CMD} ${ARG}

.PHONY: run.docker
run.docker:	## Run With Parallel Container Mode. e.g. : `make run.docker ARG="scan image"` .
	docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' -v `pwd`:/tool/data registry.veinmind.tech/veinmind/${APP} ${ARG}

##@ Clean
.PHONY: clean
clean:  ## Remove Scanner Data DB && Stop Clamd
	rm -rf data/data.db
	kill -9 `lsof -i tcp:3310 | grep 3310 | awk '{print $2}'`

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-\\.]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
