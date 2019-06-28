SERVICE = cla

generate: swagger

setup:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install
	go get -u github.com/amacneil/dbmate
	npm install serverless

clean:
	rm -rf ./gen ./bin

qc:
	gometalinter ./... --vendor

swagger: clean
	mkdir gen
	swagger -q generate server -t gen -f swagger/cla.yaml --exclude-main -A cla -P user.CLAUser

up:
	dbmate up

test:
	go test `go list ./... | grep -v node_modules`

run:
	go run main.go

build:
	dep ensure
	go build -o $(SERVICE) \
		-ldflags " \
			-X github.com/communitybridge/easy-cla/cla-backend-go/cmd.BuildStamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` \
			-X github.com/communitybridge/easy-cla/cla-backend-go/cmd.GitHash=`git rev-parse HEAD`" \
		.

buildlinux:
	dep ensure
	env GOOS=linux go build -o bin/$(SERVICE) -a -ldflags "-s -w -extldflags '-static' -X cmd.BuildStamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X cmd.GitHash=`git rev-parse HEAD`" .

build_aws_lambda:
	dep ensure
	env GOOS=linux GOARCH=amd64 go build -tags aws_lambda -ldflags="-s -w" -o backend_aws_lambda

.PHONY: generate setup clean qc swagger up test run build buildlinux build_aws_lambda
