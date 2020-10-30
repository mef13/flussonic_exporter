PROJECTFULLNAME=github.com/mef13/flussonic_exporter
export PATH := $(PATH):/usr/local/go/bin
VERSION=$(shell git describe --tags)
BUILD=$(shell git rev-parse --short HEAD)
PROJECTNAME=$(shell basename "$(PWD)")
APPPATH=/apps/$(PROJECTNAME)

# Go related variables.
GOBASE=$(shell pwd)
GOPATH=$(GOBASE)/.gocache:$(GOBASE)
GOBIN=$(GOBASE)/.build
GOFILES=$(wildcard *.go)
# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-X=main.version=$(VERSION) -X=main.commit=$(BUILD)"

# Redirect error output to a file, so we can show it in development mode.
STDERR := /tmp/.$(PROJECTNAME)-stderr.txt

# PID file will keep the process id of the server
PID := /tmp/.$(PROJECTNAME).pid

all: help

.PHONY: download
download: ## download missing dependencies. Runs `go mod download` internally
	@echo "  >  Download modules..."
	GO111MODULE=on GOPATH=$(GOPATH) GOBIN=$(GOBIN) go mod download

.PHONY: build
build: download go-build ## compile the binary in ./.build/ for current architecture

.PHONY: upgrade
upgrade: copy-app ## copy app from ./build/ to /apps/{app_name}
	@echo "  >   app upgraded."

.PHONY: fresh-install
fresh-install: add-user copy-app copy-contrib## create user, copy app and contrib(only for new install)
	@echo "  >   Configure /etc/$(PROJECTNAME)/settings.yaml before running"

.PHONY: clean
clean: ## clean dependencies
	@echo "  >   Clean $(GOBASE)/.gocache/pkg"
	@rm -rf $(GOBASE)/.gocache/pkg/*

.PHONY: go-version
go-version: ## show go version
	@go version

.PHONY: go-build
go-build:
	@echo "  >  Building binary..."
	GO111MODULE=on GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) main.go

.PHONY: copy-app
copy-app:
	@chmod +x contrib/macros.sh
	@contrib/macros.sh copy_app $(GOBIN)/$(PROJECTNAME) /usr/sbin/$(PROJECTNAME) $(PROJECTNAME)

.PHONY: copy-contrib
copy-contrib:
	@chmod +x contrib/macros.sh
	@contrib/macros.sh prepare_folder /etc/$(PROJECTNAME) $(PROJECTNAME)
	@contrib/macros.sh cp_mod $(GOBASE)/contrib/settings.yaml /etc/$(PROJECTNAME)/settings.yaml $(PROJECTNAME)
	@contrib/macros.sh prepare_folder /var/log/$(PROJECTNAME) $(PROJECTNAME)
	@cp -n $(GOBASE)/contrib/flussonic_exporter.service /lib/systemd/system/flussonic_exporter.service

.PHONY: add-user
add-user:
	@chmod +x contrib/macros.sh
	@contrib/macros.sh add_user $(PROJECTNAME)

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
