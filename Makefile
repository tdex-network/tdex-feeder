.PHONY: build-arm build-linux build-mac clean cov fmt help vet test

## build-arm: build binary for ARM
build:
	chmod u+x ./scripts/build
	./scripts/build

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
run: clean
	export FEEDER_CONFIG_PATH=./config.example.json; \
	export FEEDER_LOG_LEVEL=5; \
	go run cmd/feederd/main.go

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