.PHONY: build-arm build-linux build-mac clean cov fmt help vet test

## build-arm: build binary for ARM
build-arm:
	export GO111MODULE=on
	chmod u+x ./scripts/build
	./scripts/build linux arm

## build-linux: build binary for Linux
build-linux:
	export GO111MODULE=on
	chmod u+x ./scripts/build
	./scripts/build linux amd64

## build-mac: build binary for Mac
build-mac:
	export GO111MODULE=on
	chmod u+x ./scripts/build
	./scripts/build darwin amd64

## clean: cleans the binary
clean:
	@echo "Cleaning..."
	@go clean

## cov: generates coverage report
cov:
	@echo "Coverage..."
	go test -cover ./...

## fmt: Go Format
fmt:
	@echo "Gofmt..."
	@if [ -n "$(gofmt -l .)" ]; then echo "Go code is not formatted"; exit 1; fi


## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## run-linux: Run locally with default configuration
run-linux: clean build-linux
	export FEEDER_LOG_LEVEL=5; \
	./build/feederd-linux-amd64

## run-mac: Run locally with default configuration
run-mac: clean build-mac
	export FEEDER_LOG_LEVEL=5; \
	./build/feederd-darwin-amd64

## vet: code analysis
vet:
	@echo "Vet..."
	@go vet ./...

## test: runs go unit test with default values
test: fmt shorttest

shorttest: 
	@echo "Testing..."
	go test -v -count=1 -race -short ./...

integrationtest:
	export FEEDER_CONFIG_PATH="./config.test.json"; \
	export FEEDER_LOG_LEVEL=5; \
	go test -v -count=1 ./cmd/feederd