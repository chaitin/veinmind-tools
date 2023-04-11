.DEFAULT_GOAL := help

SHELL = /bin/bash
APP = veinmind-weakpass
CMD = cmd/*
ARG = ""

# build params
IMAGE_TAG = latest
CI_GOOS = linux
CI_GOARCH=$(shell uname -m)
TAGS = dynamic

##@ Init
.PHONY: deps
deps: 	## Install Dependencies.
ifeq ($(strip $(GITHUB_ACTIONS)),)
		go env -w GOPROXY=https://goproxy.cn,direct
endif
	go mod tidy

.PHONY: arm
arm: ## Install arm build dependency.
	apt-get update && sudo apt-get install -y gcc-aarch64-linux-gnu g++-aarch64-linux-gnu  libssl-dev

##@ Build
.PHONY: build
build: deps   ## Build Apps.
	go build -ldflags '-s -w' -v -trimpath  -a -tags ${TAGS} -o ${APP} ${CMD}

.PHONY: build.platform
build.platform: deps ## Build Apps With Platform.
	# needs libcrypt
	echo 'deb http://deb.debian.org/debian sid main' > /etc/apt/sources.list
	apt-get update
	apt-get install -y apt || true
	cd $(mktemp -d) && apt -y download libcrypt1 && dpkg-deb -x libcrypt1_*.deb . && cp -ra lib/* /lib/
	find /lib/*/libcrypt.* -ls
	apt -y --fix-broken install && apt-get install -y apt && apt-get install -y libcrypt1 libcrypt-dev
	apt-get install -y util-linux || true
	apt -y download util-linux && dpkg -i --force-overwrite util-linux* || true && apt-get install -y util-linux libc6-dev
	export CGO_ENABLED=1 CGO_LDFLAGS_ALLOW='-Wl,.*' GOOS="${CI_GOOS}" GOARCH="${CI_GOARCH}" TAGS=${TAGS} ; \
	go build -ldflags '-s -w' -v -trimpath  -a -tags ${TAGS} -o ${APP}_${CI_GOOS}_${CI_GOARCH} ${CMD}

.PHONY: build.arm ## Build Apps Static.
build.arm: deps arm
	export CC="aarch64-linux-gnu-gcc" CXX="aarch64-linux-gnu-g++"; \
	export PKG_CONFIG_PATH=/usr/lib/aarch64-linux-gnu/pkgconfig; \
	export CGO_ENABLED=1 CGO_LDFLAGS_ALLOW='.*' ; GOOS="${CI_GOOS}" GOARCH="${CI_GOARCH}"; \
	go build -ldflags '-s -w' -v -trimpath -a -tags ${TAGS} -o ${APP}_${CI_GOOS}_${CI_GOARCH} ${CMD}

.PHONY: build.docker
build.docker: ## Build Apps Docker Images.
	docker build -t ${APP}:${IMAGE_TAG} .

##@ Run
.PHONY: run
run: deps ## Run Apps. 			e.g. : `make run ARG="scan image"` .
	go run ${CMD} ${ARG}

.PHONY: run.docker
run.docker:	## Run With Parallel Container Mode. e.g. : `make run.docker ARG="scan image"` .
	docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' -v `pwd`:/tool/data registry.veinmind.tech/veinmind/${APP} ${ARG}

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-\\.]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: test
test:
	echo ${CFLAGS} > 1.txt