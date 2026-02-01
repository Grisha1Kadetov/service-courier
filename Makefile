APP_NAME=courier-service
MAIN_DIR=courier
BUILD_DIR=bin

.PHONY: ci lint test build build-check deps

deps:
	go mod download

lint:
	golangci-lint run ./...

test:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

build-check:
	go build ./...

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/$(MAIN_DIR)/

ci: deps lint test build-check
