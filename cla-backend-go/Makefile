# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT
SERVICE = cla
LAMBDA_BIN = backend-aws-lambda
METRICS_BIN = metrics-aws-lambda
FUNCTIONAL_TESTS_BIN = functional-tests
BUILD_TIME=`date +%FT%T%z`
VERSION := $(shell sh -c 'git describe --always --tags')
BRANCH := $(shell sh -c 'git rev-parse --abbrev-ref HEAD')
COMMIT := $(shell sh -c 'git rev-parse --short HEAD')
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.branch=$(BRANCH) -X main.buildDate=$(BUILD_TIME)"
BUILD_TAGS=-tags aws_lambda

LINT_TOOL=$(shell go env GOPATH)/bin/golangci-lint
LINT_VERSION=v1.25.0
GO_PKGS=$(shell go list ./... | grep -v /vendor/ | grep -v /node_modules/)
GO_FILES=$(shell find . -type f -name '*.go' -not -path './vendor/*')
TEST_ENV=AWS_REGION=us-east-1 DYNAMODB_AWS_REGION=us-east-1 AWS_PROFILE=bar AWS_ACCESS_KEY_ID=foo AWS_SECRET_ACCESS_KEY=bar

.PHONY: generate setup setup-dev setup-deploy clean-all clean swagger up fmt test run deps build build-mac build-aws-lambda qc lint

all: clean swagger deps fmt build-mac test lint

generate: swagger

setup: $(LINT_TOOL) setup-dev setup-deploy

setup_dev: setup-dev
setup-dev:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	go get -u golang.org/x/tools/cmd/goimports
	go get golang.org/x/tools/cmd/cover
	npm install -g multi-file-swagger


setup_deploy: setup-deploy
setup-deploy:
	npm install serverless

clean:
	@rm -rf bin/ cla cla-mac cla-linux \
		./v2/project-service/client ./v2/project-service/models \
		./v2/organization-service/client ./v2/organization-service/models \
		./v2/user-service/client ./v2/user-service/models

clean-swagger:
	@rm -rf gen/

clean-all: clean clean-swagger
	@rm -rf vendor/

swagger: clean-swagger swagger-build swagger-validate

swagger-prep:
	@mkdir gen

swagger-build-v1-services:
	@echo "Generating v1 legacy API models..."
	cd swagger; multi-file-swagger -o yaml cla.yaml > cla.compiled.yaml
	swagger -q generate server \
		-t gen \
		-f swagger/cla.compiled.yaml \
		--copyright-file=copyright-header.txt \
		--exclude-main \
		-A cla \
		-P user.CLAUser

swagger-build-v2-services:
	@echo "Generating v2 API models..."
	cd swagger; multi-file-swagger -o yaml cla.v2.yaml > cla.v2.compiled.yaml
	swagger -q generate server \
		-t gen \
		-f swagger/cla.v2.compiled.yaml \
		--copyright-file=copyright-header.txt \
		--server-package=v2/restapi \
		--model-package=v2/models \
		--exclude-main \
		-A easycla \
		-P auth.User

swagger-build-project-service:
	@echo "Generating swagger client for the project-service..."
	@mkdir -p v2/project-service
	curl -sfL https://api-gw.dev.platform.linuxfoundation.org/project-service/swagger.json --output swagger/project-service.yaml
	sed -i.bak 's/X-ACL/Empty-Header/g' swagger/project-service.yaml
	swagger -q generate client \
		--copyright-file=copyright-header.txt \
		-t v2/project-service \
		-f swagger/project-service.yaml

swagger-build-organization-service:
	@echo "Generating swagger client for the organization-service..."
	@mkdir -p v2/organization-service
	curl -sfL https://api-gw.dev.platform.linuxfoundation.org/organization-service/swagger.json --output swagger/organization-service.yaml
	sed -i.bak 's/X-ACL/Empty-Header/g' swagger/organization-service.yaml
	swagger -q generate client \
		--copyright-file=copyright-header.txt \
		-t v2/organization-service \
		-f swagger/organization-service.yaml

swagger-build-user-service:
	@echo "Generating swagger client for the user-service..."
	@mkdir -p v2/user-service
	curl -sfL https://api-gw.dev.platform.linuxfoundation.org/user-service/swagger.json --output swagger/user-service.yaml
	sed -i.bak 's/X-ACL/Empty-Header/g' swagger/user-service.yaml
	swagger -q generate client \
		--copyright-file=copyright-header.txt \
		-t v2/user-service \
		-f swagger/user-service.yaml

swagger-build-acs-service:
	@echo "Generating swagger client for the acs-service..."
	@mkdir -p v2/acs-service
	curl -sfL https://api-gw.dev.platform.linuxfoundation.org/acs/v1/api-docs/swagger/swagger.json --output swagger/acs-service.yaml
	sed -i.bak 's/X-ACL/X-API-KEY/g' swagger/acs-service.yaml
	swagger -q generate client \
		--copyright-file=copyright-header.txt \
		-t v2/acs-service \
		-f swagger/acs-service.yaml

swagger-build: clean-swagger swagger-prep swagger-build-v1-services swagger-build-v2-services swagger-build-project-service swagger-build-organization-service swagger-build-user-service swagger-build-acs-service

swagger-validate:
	@echo "Validating v1 legacy API specification..."
	swagger validate swagger/cla.compiled.yaml
	@echo "Validating v2 API specification..."
	swagger validate swagger/cla.v2.compiled.yaml

fmt:
	@gofmt -w -l -s $(GO_FILES)
	@goimports -w -l $(GO_FILES)

test:
	@ $(TEST_ENV) go test -v $(shell go list ./... | grep -v /vendor/ | grep -v /node_modules/) -coverprofile=cover.out

run:
	go run main.go

deps:
	go env -w GOPRIVATE=github.com/LF-Engineering/*
	go mod download

build: deps
	@echo "Building Linux amd64 binary..."
	env GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(SERVICE) main.go
	@chmod +x $(SERVICE)

build-mac: deps
	@echo "Building Mac OSX amd64 binary..."
	env GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(SERVICE)-mac main.go
	@chmod +x $(SERVICE)-mac

rebuild-mac: clean fmt build-mac lint
	./$(SERVICE)-mac

build-aws-lambda: deps
	@echo "Building a statically linked Linux amd64 binary..."
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) $(BUILD_TAGS) -o $(LAMBDA_BIN) main.go
	@chmod +x $(LAMBDA_BIN)

build-metrics-lambda: deps
	@echo "Building a statically linked Linux amd64 binary..."
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(METRICS_BIN) cmd/metrics_lambda/main.go
	@chmod +x $(METRICS_BIN)

build-metrics-lambda-mac: deps
	@echo "Building a statically linked Mac OSX amd64 binary..."
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(METRICS_BIN)-mac cmd/metrics_lambda/main.go
	@chmod +x $(METRICS_BIN)-mac

build-functional-tests-linux: deps
	@echo "Building Functional Tests for Linux amd64 binary..."
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(FUNCTIONAL_TESTS_BIN) cmd/functional_tests/main.go
	@chmod +x $(FUNCTIONAL_TESTS_BIN)

build-functional-tests-mac: deps
	@echo "Building Functional Tests for OSX amd64 binary..."
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(FUNCTIONAL_TESTS_BIN)-mac cmd/functional_tests/main.go
	@chmod +x $(FUNCTIONAL_TESTS_BIN)-mac

$(LINT_TOOL):
	@echo "Downloading golangci-lint version $(LINT_VERSION)..."
	# Latest releases: https://github.com/golangci/golangci-lint/releases
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(shell go env GOPATH)/bin $(LINT_VERSION)

lint: $(LINT_TOOL)
	$(LINT_TOOL) run --config=.golangci.yaml ./...

