# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT
MAKEFILE_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
BUILD_TIME=`date +%FT%T%z`
VERSION := $(shell sh -c 'git describe --always --tags')
BRANCH := $(shell sh -c 'git rev-parse --abbrev-ref HEAD')
COMMIT := $(shell sh -c 'git rev-parse --short HEAD')
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.branch=$(BRANCH) -X main.buildDate=$(BUILD_TIME)"

LINT_TOOL=$(shell go env GOPATH)/bin/golangci-lint
LINT_VERSION=v1.29.0
SWAGGER_TOOL_VERSION=v0.24.0
GO_PKGS=$(shell go list ./... | grep -v /vendor/ | grep -v /node_modules/)
GO_FILES=$(shell find . -type f -name '*.go' -not -path './vendor/*')

.PHONY: generate setup tool-setup setup-dev setup-deploy clean-all clean swagger up fmt test run deps build build-mac build-aws-lambda user-subscribe-lambda qc lint

setup: $(LINT_TOOL) setup-dev

tool-setup:
	@echo "Installing gobin for installing tools..."
	@# gobin is the equivalent of 'go get' whilst in module-aware mode but this does not modify your go.mod
	GO111MODULE=off go get -u github.com/myitcv/gobin

setup-dev: tool-setup
	@echo "Installing goimports..."
	gobin golang.org/x/tools/cmd/goimports
	@echo "Installing cover..."
	gobin golang.org/x/tools/cmd/cover

clean:
	@rm -rf main

fmt:
	@echo "Formatting code and optimizing imports..."
	@gofmt -w -l -s $(GO_FILES)
	@goimports -w -l $(GO_FILES)

test:
	@echo "Running unit tests..."
	@ go test -v $(shell go list ./... | grep -v /vendor/ | grep -v /node_modules/) -coverprofile=cover.out

run:
	go run main.go

deps:
	go mod download

build: build-linux
build-linux: deps
	@echo "Building Linux amd64 binary..."
	env GOOS=linux GOARCH=amd64 go build $(LDFLAGS) main.go

build-mac: deps
	@echo "Building Mac OSX amd64 binary..."
	env GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) main.go

$(LINT_TOOL):
	@echo "Downloading golangci-lint version $(LINT_VERSION)..."
	@# Latest releases: https://github.com/golangci/golangci-lint/releases
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(shell go env GOPATH)/bin $(LINT_VERSION)

lint: $(LINT_TOOL)
	@cd $(MAKEFILE_DIR) && echo "Running lint..." && $(LINT_TOOL) run --allow-parallel-runners --config=.golangci.yaml ./... && echo "Lint check passed."
	@cd $(MAKEFILE_DIR) && ./check-headers.sh

