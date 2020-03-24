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
LINT_VERSION=v1.23.7
GO_PKGS=$(shell go list ./... | grep -v /vendor/ | grep -v /node_modules/)
GO_FILES=$(shell find . -type f -name '*.go' -not -path './vendor/*')
TEST_ENV=AWS_REGION=us-east-1 DYNAMODB_AWS_REGION=us-east-1 AWS_PROFILE=bar AWS_ACCESS_KEY_ID=foo AWS_SECRET_ACCESS_KEY=bar

.PHONY: generate setup setup_dev setup_deploy clean-all clean swagger up fmt test run deps build build-mac build-aws-lambda qc lint

generate: swagger

setup: $(LINT_TOOL) setup_dev setup_deploy

setup_dev:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/golang/dep/cmd/dep	
	go get golang.org/x/tools/cmd/cover
	go get -u github.com/stripe/safesql

setup_deploy:
	npm install serverless

clean:
	rm -rf ./gen ./bin

clean-all: clean
	rm -rf vendor

swagger: clean
	@mkdir gen
	@echo "Genering v1 legacy API models..."
	swagger -q generate server \
		-t gen \
		-f swagger/cla.yaml \
		--copyright-file=copyright-header.txt \
		--exclude-main \
		-A cla \
		-P user.CLAUser
	@echo "Genering v2 API models..."
	swagger -q generate server \
		-t gen \
		-f swagger/cla.v2.yaml \
		--copyright-file=copyright-header.txt \
		--server-package=v2/restapi \
		--model-package=v2/models \
		--exclude-main \
		-A easycla \
		-P auth.User

swagger-validate:
	@echo "Validating v1 legacy API models..."
	swagger validate swagger/cla.yaml
	@echo "Validating v2 API models..."
	swagger validate swagger/cla.v2.yaml

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
	@echo "Building Linux amd64 binary..."
	env GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(SERVICE) main.go
	@chmod +x cla

build-mac: deps
	@echo "Building Mac OSX amd64 binary..."
	env GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(SERVICE) main.go
	@chmod +x cla

rebuild-mac:
	rm -rf ./cla
	make fmt build-mac lint
	./cla

build-aws-lambda: deps
	@echo "Building a staticlly linked Linux amd64 binary..."
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) $(BUILD_TAGS) -o backend-aws-lambda main.go
	chmod +x backend-aws-lambda

build-metrics-lambda: deps
	@echo "Building a staticlly linked Linux amd64 binary..."
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o metrics-aws-lambda cmd/metrics_lambda/main.go
	chmod +x metrics-aws-lambda

$(LINT_TOOL):
	@echo "Downloading golangci-lint version $(LINT_VERSION)..."
	# Latest releases: https://github.com/golangci/golangci-lint/releases
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(shell go env GOPATH)/bin $(LINT_VERSION)

lint: $(LINT_TOOL)
	$(LINT_TOOL) run --config=.golangci.yaml ./...

all: clean swagger swagger-validate deps fmt build-mac test lint
