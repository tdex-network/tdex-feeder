.PHONY: build run clean cov fmt help vet test shorttest integrationtest

## proto: compile proto stubs
proto:
	@echo "Compiling stubs..."
	@buf generate buf.build/tdex-network/tdex-protobuf


## build: builds binary for all platforms
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

## fmt: checks if code is well formatted
fmt:
	@echo "Gofmt..."
	@if [ -n "$(gofmt -l .)" ]; then echo "Go code is not formatted"; exit 1; fi


## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# run: runs locally without building binary
run: clean
	export FEEDER_LOG_LEVEL=5; \
	go run cmd/feederd/main.go

## vet: code analysis
vet:
	@echo "Vetting..."
	@go vet ./...

## test: runs unit tests
test: fmt
	@echo "Testing..."
	go test -v -count=1 -race -short ./...

## integrationtest: run integration tests
integrationtest:
	export FEEDER_CONFIG_PATH="./config.test.json"; \
	export FEEDER_LOG_LEVEL=5; \
	go test -v -count=1 ./cmd/feederd