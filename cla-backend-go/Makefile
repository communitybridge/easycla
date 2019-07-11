# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT
SERVICE = cla
LINT_TOOL=$(shell go env GOPATH)/bin/golangci-lint
GO_PKGS=$(shell go list ./... | grep -v /vendor/ | grep -v /node_modules/)
GO_FILES=$(shell find . -type f -name '*.go' -not -path './vendor/*')

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

up:
	dbmate up

fmt:
	@go fmt $(GO_PKGS)
	@goimports -w -l $(GO_FILES)

test:
	@go test -v $(shell go list ./... | grep -v /vendor/ | grep -v /node_modules/) -coverprofile=cover.out

run:
	go run main.go

deps:
	dep ensure -v

build: deps
	go build -o $(SERVICE) \
		-ldflags " \
			-X github.com/communitybridge/easycla/cla-backend-go/cmd.BuildStamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` \
			-X github.com/communitybridge/easycla/cla-backend-go/cmd.GitHash=`git rev-parse HEAD`" \
		.

buildlinux: deps
	env GOOS=linux go build -o bin/$(SERVICE) -a -ldflags "-s -w -extldflags '-static' -X cmd.BuildStamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X cmd.GitHash=`git rev-parse HEAD`" .

build_aws_lambda: deps
	env GOOS=linux GOARCH=amd64 go build -tags aws_lambda -ldflags="-s -w" -o backend_aws_lambda

$(LINT_TOOL):
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.17.1

qc: $(LINT_TOOL)
	$(LINT_TOOL) run --config=.golangci.yaml ./...

lint: qc

.PHONY: generate setup clean swagger up test run build buildlinux build_aws_lambda qc lint
