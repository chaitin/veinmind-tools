.DEFAULT_GOAL := help
SHELL = /bin/bash
# version control
VERSION ?= "latest"
# plugin int params
LANG = go
NAME ?= ""
PUB = no
# color
DEFAULT=\033[0m
GREEN=\033[32m
RED=\033[31m
BLUE=\033[96m

# platform
CI_GOOS=linux
CI_GOARCH=$(shell uname -m)
TAGS ?=

ifeq ("$(shell uname)", "Darwin")
define sed
	@sed -ri "" $(1) $(2)
endef
else
define sed
	@sed -ri $(1) $(2)
endef
endif
define update
	$(call sed, $(1) $(shell grep $(2) -rl ./plugins ./example ./veinmind-runner ./.github))
endef
define init_plugin
	@echo -e "$(GREEN)~~~~~~~~~~~~~~~~~~~~~~~~~~~ Welcome to Veinmind Tools ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"; \
	echo -e "enter language(go/python): $(LANG)"; \
	echo -e "enter plugins name: $(NAME) "; \
	echo -e "enter is publish(PUB: DEFAULT no): $(PUB)"; \
	mkdir $(1); \
	echo -e "$(BLUE)init Veinmind GO Plugin $(NAME) at: plugins/$(LANG)/veinmind-$(NAME)";\
	echo -e "# veinmind-$(NAME)  \n\n这是描述文件" >$(1)/README.md;\
	echo -e "# veinmind-$(NAME)  \n\nthis is description file" > $(1)/README.en.md;
endef

##@ Init
.PHONY: install
install: ## 			install libVeinMind
	@echo 'deb [trusted=yes] https://download.veinmind.tech/libveinmind/apt/ ./' | sudo tee /etc/apt/sources.list.d/libveinmind.list; \
    apt-get update;\
    apt-get install -y libveinmind-dev

.PHONY: plugin
plugin: ## 			init a new Plugins
ifeq ($(LANG), go)
	$(call init_plugin, plugins/go/veinmind-$(NAME))
	@cp -r ./example/go/* plugins/go/veinmind-$(NAME)
	$(call sed, "s/veinmind-example/veinmind-$(NAME)/g", plugins/go/veinmind-$(NAME)/Dockerfile)
	$(call sed, "s/veinmind-example/veinmind-$(NAME)/g", plugins/go/veinmind-$(NAME)/Makefile)
	$(call sed, "s/veinmind-example/veinmind-$(NAME)/g", plugins/go/veinmind-$(NAME)/go.mod)
	$(call sed, "s/veinmind-example/veinmind-$(NAME)/g", plugins/go/veinmind-$(NAME)/cmd/cli.go)
else ifeq ($(LANG), python)
	$(call init_plugin, plugins/python/veinmind-$(NAME))
	@cp -r ./example/python/* plugins/python/veinmind-$(NAME)
	$(call sed, "s/veinmind-example/veinmind-$(NAME)/g", plugins/python/veinmind-$(NAME)/Dockerfile)
	$(call sed, "s/veinmind-example/veinmind-$(NAME)/g", plugins/python/veinmind-$(NAME)/scan.py)
endif

##@ Update
.PHONY: libveinmind
libveinmind: 	##			upgrade libVeinMind
ifeq ($(VERSION), "")
	@echo "VERSION is empty, use 'make update.libveinmind VERSION=x.x.x'"
else
	$(call update, "s/github\.com\/chaitin\/libveinmind v[0-9]\.[0-9]\.[0-9]$$/github.com\/chaitin\/libveinmind v$(VERSION)/g", "github.com/chaitin/libveinmind v[0-9]\.[0-9]\.[0-9]$$");
	$(call update, "s/veinmind==[0-9]\.[0-9]\.[0-9]/veinmind==$(VERSION)/g", "veinmind==")
endif

.PHONY: libveinmind-docker
libveinmind-docker:  ##  		upgrade libVeinMind in Dockerfile
ifeq ($(VERSION), "")
	@echo "VERSION is empty, use 'make update.libveinmind-docker VERSION=x.x.x'"
else
	$(call update, "s/veinmind\/python3[0-9\.]*:[0-9]\.[0-9]\.[0-9]*/veinmind\/python3.6:$(VERSION)/g", "veinmind\/python3[0-9\.]*:[0-9]\.[0-9]\.[0-9]")
	$(call update, "s/veinmind\/go1.*:[0-9]\.[0-9]\.[0-9]*/veinmind\/go1.18:$(VERSION)/g", "veinmind\/go1.*:[0-9]\.[0-9]\.[0-9]")
	$(call update, "s/veinmind\/base:[0-9]\.[0-9]\.[0-9]*/veinmind\/base:$(VERSION)/g", "veinmind\/base:[0-9]\.[0-9]\.[0-9]")
endif

.PHONY: veinmind-common-go
veinmind-common-go: ## 		upgrade veinmind-common-go
ifeq ($(VERSION), "")
	@echo "VERSION is empty, use 'make update.veinmind-common-go VERSION=x.x.x'"
else
	$(call update, "s/github\.com\/chaitin\/veinmind-common-go v[0-9]\.[0-9]\.[0-9](-r[0-9])?$$/github.com\/chaitin\/veinmind-common-go v$(VERSION)/g", "github\.com\/chaitin\/veinmind-common-go v[0-9]\.[0-9]\.[0-9]\(-r[0-9]\)\?$$");
endif

.PHONY: veinmind-common-python
veinmind-common-python:  ##	upgrade veinmind-common-python
ifeq ($(LIBVEINMIND_COMMON_PYTHON_VERSION), "")
	@echo "VERSION is empty, use 'make update.veinmind-common-python VERSION=x.x.x`"
else
	$(call update, "s/veinmind-common==[0-9]\.[0-9]\.[0-9](\.post[0-9])?/veinmind-common==$(VERSION)/g", "veinmind-common==")
endif


##@ Build
all: ## 			build all plugins
	$(MAKE) $(shell ls plugins/go/)

veinmind-%: ##			build go plugins. e.g. `make veinmind-basic`
	$(MAKE) -C plugins/go/$@ CGO_ENABLED=1 build

platform.veinmind-runner:
	$(MAKE) -C veinmind-runner build.platform CI_GOOS=${CI_GOOS} CI_GOARCH=${CI_GOARCH} TAGS=${TAGS}

platform.veinmind-%: ##   	build go plugins with platform. e.g. `make veinmind-basic CI_GOOS=linux CI_GOARCH=amd64 TAGS=xxxx`
	$(MAKE) -C plugins/go/$(subst platform.,,$@) build.platform CI_GOOS=${CI_GOOS} CI_GOARCH=${CI_GOARCH} TAGS=${TAGS}

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z_\-\\.%]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
