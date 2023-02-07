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
	echo -e "$(BLUE)init Veinmind GO Plugin $(NAME) at: plugins/$(LANG)/$(NAME)";\
	cp ./example/parallel-container-run.sh ./example/README.md ./example/README.en.md $(1);\
	echo -e "# veinmind-$(NAME)  \n\n这是描述文件" >$(1)/README.md;\
	echo -e "# veinmind-$(NAME)  \n\nthis is description file" > $(1)/README.en.md;
	$(call sed, "s/veinmind-example/veinmind-$(NAME)/g", $(1)/parallel-container-run.sh);
endef


.PHONY: all
# build all plugins
all: build.veinmind-basic

.PHONY: install
# install LibVeinMind
install:
	@echo 'deb [trusted=yes] https://download.veinmind.tech/libveinmind/apt/ ./' | sudo tee /etc/apt/sources.list.d/libveinmind.list; \
    sudo apt-get update;\
    sudo apt-get install -y libveinmind-dev

.PHONY: update.libveinmind
# update LibVeinMind
update.libveinmind:
ifeq ($(VERSION), "")
	@echo "VERSION is empty, use 'make update.libveinmind VERSION=x.x.x'"
else
	$(call update, "s/github\.com\/chaitin\/libveinmind v[0-9]\.[0-9]\.[0-9]$$/github.com\/chaitin\/libveinmind v$(VERSION)/g", "github.com/chaitin/libveinmind v[0-9]\.[0-9]\.[0-9]$$");
	$(call update, "s/veinmind==[0-9]\.[0-9]\.[0-9]/veinmind==$(VERSION)/g", "veinmind==")
endif

.PHONY: update.libveinmind-docker
# update LibVeinMind-Dockerfile
update.libveinmind-docker:
ifeq ($(VERSION), "")
	@echo "VERSION is empty, use 'make update.libveinmind-docker VERSION=x.x.x'"
else
	$(call update, "s/veinmind\/python3[0-9\.]*:[0-9]\.[0-9]\.[0-9]/veinmind\/python3.6:$(VERSION)/g", "veinmind\/python3[0-9\.]*:[0-9]\.[0-9]\.[0-9]")
	$(call update, "s/veinmind\/go1.*:[0-9]\.[0-9]\.[0-9]/veinmind\/go1.18:$(VERSION)/g", "veinmind\/go1.*:[0-9]\.[0-9]\.[0-9]")
	$(call update, "s/veinmind\/base:[0-9]\.[0-9]\.[0-9]/veinmind\/base:$(VERSION)/g", "veinmind\/base:[0-9]\.[0-9]\.[0-9]")
endif

.PHONY: update.veinmind-common-go
# upgrade veinmind-common-go
update.veinmind-common-go:
ifeq ($(VERSION), "")
	@echo "VERSION is empty, use 'make update.veinmind-common-go VERSION=x.x.x'"
else
	$(call update, "s/github\.com\/chaitin\/veinmind-common-go v[0-9]\.[0-9]\.[0-9](-r[0-9])?$$/github.com\/chaitin\/veinmind-common-go v$(VERSION)/g", "github\.com\/chaitin\/veinmind-common-go v[0-9]\.[0-9]\.[0-9]\(-r[0-9]\)\?$$");
endif

.PHONY: update.veinmind-common-python
# upgrade veinmind-common-python
update.veinmind-common-python:
ifeq ($(LIBVEINMIND_COMMON_PYTHON_VERSION), "")
	@echo "VERSION is empty, use 'make update.veinmind-common-python VERSION=x.x.x`"
else
	$(call update, "s/veinmind-common==[0-9]\.[0-9]\.[0-9](\.post[0-9])?/veinmind-common==$(VERSION)/g", "veinmind-common==")
endif

.PHONY: plugin.init
# init plugin
plugin.init:
ifeq ($(LANG), go)
	$(call init_plugin, plugins/go/veinmind-$(NAME))
	@cp -r ./example/go/* plugins/go/veinmind-$(NAME)
	$(call sed, "s/veinmind-example/veinmind-$(NAME)/g", plugins/go/veinmind-$(NAME)/Dockerfile)
	$(call sed, "s/veinmind-example/veinmind-$(NAME)/g", plugins/go/veinmind-$(NAME)/script/build_amd64.sh)
	$(call sed, "s/veinmind-example/veinmind-$(NAME)/g", plugins/go/veinmind-$(NAME)/script/build.sh)
	$(call sed, "s/veinmind-example/veinmind-$(NAME)/g", plugins/go/veinmind-$(NAME)/go.mod)
	$(call sed, "s/veinmind-example/veinmind-$(NAME)/g", plugins/go/veinmind-$(NAME)/cmd/cli.go)
else ifeq ($(LANG), python)
	$(call init_plugin, plugins/python/veinmind-$(NAME))
	@cp -r ./example/python/* plugins/python/veinmind-$(NAME)
	$(call sed, "s/veinmind-example/veinmind-$(NAME)/g", plugins/python/veinmind-$(NAME)/Dockerfile)
	$(call sed, "s/veinmind-example/veinmind-$(NAME)/g", plugins/python/veinmind-$(NAME)/scan.py)
endif

.PHONY: build.veinmind-basic
build.veinmind-basic:
	$(MAKE) -C plugins/go/veinmind-basic build
