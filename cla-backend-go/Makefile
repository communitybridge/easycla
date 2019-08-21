# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT
SERVICE = cla
BUILD_TIME=`date +%FT%T%z`
VERSION := $(shell sh -c 'git describe --always --tags')
BRANCH := $(shell sh -c 'git rev-parse --abbrev-ref HEAD')
COMMIT := $(shell sh -c 'git rev-parse --short HEAD')
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.branch=$(BRANCH) -X main.buildDate=$(BUILD_TIME)"
BUILD_TAGS=-tags aws_lambda

LINT_TOOL=$(shell go env GOPATH)/bin/golangci-lint
GO_PKGS=$(shell go list ./... | grep -v /vendor/ | grep -v /node_modules/)
GO_FILES=$(shell find . -type f -name '*.go' -not -path './vendor/*')
TEST_ENV=AWS_REGION=us-east-1 AWS_ACCESS_KEY_ID=foo AWS_SECRET_ACCESS_KEY=bar

.PHONY: generate setup setup_dev setup_deploy clean swagger up fmt test run deps build build-mac build_aws_lambda qc lint

generate: swagger

setup: $(LINT_TOOL) setup_dev setup_deploy

setup_dev:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/golang/dep/cmd/dep	
	go get golang.org/x/tools/cmd/cover
	go get -u github.com/amacneil/dbmate
	go get -u github.com/stripe/safesql

setup_deploy:
	npm install serverless

clean:
	rm -rf ./gen ./bin

swagger: clean
	mkdir gen
	swagger -q generate server -t gen -f swagger/cla.yaml --exclude-main -A cla -P user.CLAUser

swagger-validate:
	swagger validate swagger/cla.yaml

up:
	dbmate up

fmt:
	@gofmt -w -l -s $(GO_FILES)
	@goimports -w -l $(GO_FILES)

test:
	@ $(TEST_ENV) go test -v $(shell go list ./... | grep -v /vendor/ | grep -v /node_modules/) -coverprofile=cover.out

run:
	go run main.go

deps:
	dep ensure -v

build: deps
	env GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(SERVICE) main.go
	chmod +x cla

build-mac: deps
	env GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(SERVICE) main.go
	chmod +x cla

build-aws-lambda: deps
	env GOOS=linux GOARCH=amd64 go build $(LDFLAGS) $(BUILD_TAGS) -o backend_aws_lambda main.go
	chmod +x backend_aws_lambda

$(LINT_TOOL):
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.17.1

qc: $(LINT_TOOL)
	$(LINT_TOOL) run --config=.golangci.yaml ./...

lint: qc
