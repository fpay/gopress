.PHONY: all test build lint fmt install watch

test:
	@go test -race -coverprofile=coverage.out -covermode=atomic

build:
	@go build -race

lint:
	@golint

fmt:
	@gofmt -s -w *.go

watch:
	@watchman-make -p '*.go' --make="go test" -t ""

install:
	@cd ./gopress && go install

all: test build
