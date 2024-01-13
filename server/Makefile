# Makefile
.PHONY: build build-prod test run clean

PROJECT_NAME = Dalkak
MODE ?= DEV

build:
	@echo "Building $(PROJECT_NAME) in $(MODE) mode..."
	@go build -o $(PROJECT_NAME) -tags $(MODE) -ldflags "-X main.Mode=$(MODE)" ./cmd/api/.

build-prod: MODE = PROD
build-prod:
	@echo "Building $(PROJECT_NAME) in $(MODE) mode..."
	@go build -o $(PROJECT_NAME) -tags $(MODE) -ldflags "-X main.Mode=$(MODE)" ./cmd/api/.

run:
	@echo "Running $(PROJECT_NAME)..."
	@go run -ldflags "-X main.Mode=$(MODE)" ./cmd/api/.

run-local: MODE = LOCAL
run-local:
	@echo "Running $(PROJECT_NAME)..."
	@go run -ldflags "-X main.Mode=$(MODE)" ./cmd/api/.

test:
	@echo "Testing $(PROJECT_NAME)..."
	@MODE=$(MODE) go test ./...

clean:
	@echo "Cleaning $(PROJECT_NAME)..."
	@go clean