# Makefile
.PHONY: build build-prod test run clean

PROJECT_NAME = Dalkak
MODE ?= DEV

build:
	@echo "Building $(PROJECT_NAME) in $(MODE) mode..."
	@go build -o $(PROJECT_NAME) -tags $(MODE) ./cmd/api/.

build-prod: MODE = PROD
build-prod:
	@echo "Building $(PROJECT_NAME) in $(MODE) mode..."
	@go build -o $(PROJECT_NAME) -tags $(MODE) ./cmd/api/.

run:
	@echo "Running $(PROJECT_NAME)..."
	@go run ./cmd/api/.

test:
	@echo "Testing $(PROJECT_NAME)..."
	@go test ./...

clean:
	@echo "Cleaning $(PROJECT_NAME)..."
	@go clean