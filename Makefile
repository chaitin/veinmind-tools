SHELL = /bin/bash

# version control
VERSION ?= ""

# plugin int params
Language = go
Name ?= ""
Pub = no

# color
default=\033[0m
green=\033[32m
red=\033[31m
blue=\033[96m

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
	@echo -e "$(green)~~~~~~~~~~~~~~~~~~~~~~~~~~~ Welcome to Veinmind Tools ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"; \
	echo -e "enter language(go/python): $(Language)"; \
	echo -e "enter name: $(Name) "; \
	echo -e "enter is publish(Pub: default no): $(Pub)"; \
	mkdir $(1); \
	echo -e "$(blue)init Veinmind GO Plugin $(Name) at: plugins/$(Language)/$(Name)";\
	cp ./example/parallel-container-run.sh ./example/README.md ./example/README.en.md $(1);\
	echo -e "# veinmind-$(Name)  \n\n这是描述文件" >$(1)/README.md;\
	echo -e "# veinmind-$(Name)  \n\nthis is description file" > $(1)/README.en.md;
	$(call sed, "s/veinmind-example/veinmind-$(Name)/g", $(1)/parallel-container-run.sh);
endef

define init_go_plugin
	$(call init_plugin, plugins/go/veinmind-$(Name))
	@cp -r ./example/go/* plugins/go/veinmind-$(Name)
	$(call sed, "s/veinmind-example/veinmind-$(Name)/g", plugins/go/veinmind-$(Name)/Dockerfile)
	$(call sed, "s/veinmind-example/veinmind-$(Name)/g", plugins/go/veinmind-$(Name)/script/build_amd64.sh)
	$(call sed, "s/veinmind-example/veinmind-$(Name)/g", plugins/go/veinmind-$(Name)/script/build.sh)
	$(call sed, "s/veinmind-example/veinmind-$(Name)/g", plugins/go/veinmind-$(Name)/go.mod)
	$(call sed, "s/veinmind-example/veinmind-$(Name)/g", plugins/go/veinmind-$(Name)/cmd/cli.go)
endef

define init_python_plugin
	$(call init_plugin, plugins/python/veinmind-$(Name))
	@cp -r ./example/python/* plugins/python/veinmind-$(Name)
	$(call sed, "s/veinmind-example/veinmind-$(Name)/g", plugins/python/veinmind-$(Name)/Dockerfile)
	$(call sed, "s/veinmind-example/veinmind-$(Name)/g", plugins/python/veinmind-$(Name)/scan.py)
endef

# install LibVeinMind
.PHONY: install
install:
	@echo 'deb [trusted=yes] https://download.veinmind.tech/libveinmind/apt/ ./' | sudo tee /etc/apt/sources.list.d/libveinmind.list; \
    sudo apt-get update;\
    sudo apt-get install -y libveinmind-dev

# update LibVeinMind
.PHONY: update.libveinmind
update.libveinmind:
ifeq ($(VERSION), "")
	@echo "VERSION is empty, use 'make update.libveinmind VERSION=x.x.x'"
else
	$(call update, "s/github\.com\/chaitin\/libveinmind v[0-9]\.[0-9]\.[0-9]$$/github.com\/chaitin\/libveinmind v$(VERSION)/g", "github.com/chaitin/libveinmind v[0-9]\.[0-9]\.[0-9]$$");
	$(call update, "s/veinmind==[0-9]\.[0-9]\.[0-9]/veinmind==$(VERSION)/g", "veinmind==")
endif

# update LibVeinMind-Dockerfile
.PHONY: update.libveinmind-docker
update.libveinmind-docker:
ifeq ($(VERSION), "")
	@echo "VERSION is empty, use 'make update.libveinmind-docker VERSION=x.x.x'"
else
	$(call update, "s/veinmind\/python3[0-9\.]*:[0-9]\.[0-9]\.[0-9]/veinmind\/python3.6:$(VERSION)/g", "veinmind\/python3[0-9\.]*:[0-9]\.[0-9]\.[0-9]")
	$(call update, "s/veinmind\/go1.*:[0-9]\.[0-9]\.[0-9]/veinmind\/go1.18:$(VERSION)/g", "veinmind\/go1.*:[0-9]\.[0-9]\.[0-9]")
	$(call update, "s/veinmind\/base:[0-9]\.[0-9]\.[0-9]/veinmind\/base:$(VERSION)/g", "veinmind\/base:[0-9]\.[0-9]\.[0-9]")
endif

# upgrade veinmind-common-go
.PHONY: update.veinmind-common-go
update.veinmind-common-go:
ifeq ($(VERSION), "")
	@echo "VERSION is empty, use 'make update.veinmind-common-go VERSION=x.x.x'"
else
	$(call update, "s/github\.com\/chaitin\/veinmind-common-go v[0-9]\.[0-9]\.[0-9](-r[0-9])?$$/github.com\/chaitin\/veinmind-common-go v$(VERSION)/g", "github\.com\/chaitin\/veinmind-common-go v[0-9]\.[0-9]\.[0-9]\(-r[0-9]\)\?$$");
endif

# upgrade veinmind-common-python
.PHONY: update.veinmind-common-python
update.veinmind-common-python:
ifeq ($(LIBVEINMIND_COMMON_PYTHON_VERSION), "")
	@echo "VERSION is empty, use 'make update.veinmind-common-python VERSION=x.x.x`"
else
	$(call update, "s/veinmind-common==[0-9]\.[0-9]\.[0-9](\.post[0-9])?/veinmind-common==$(VERSION)/g", "veinmind-common==")
endif

# init plugin
.PHONY: plugin.init
plugin.init:
ifeq ($(Language), go)
	$(call init_go_plugin)
else ifeq ($(Language), python)
	$(call init_python_plugin)
endif