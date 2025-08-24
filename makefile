# Makefile for Go service with environment file support

# Variables
ENV_FILE ?= .env
GO_FILES := $(shell find . -name "*.go" -type f)
MAIN_FILE := main.go
SERVICE_CMD := service

# Default target
.PHONY: all
all: help

# Help message
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make service        - Run the service (tidy dependencies, source env file, run app)"
	@echo "  make tidy           - Tidy up the go.mod file"
	@echo "  make run            - Run the application"
	@echo "  make build          - Build the application"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build artifacts"
	@echo ""
	@echo "Environment variables:"
	@echo "  ENV_FILE            - Path to environment file (default: .env)"
	@echo ""
	@echo "Example:"
	@echo "  make service ENV_FILE=.env.development"

# Combined service command (tidy + env + run)
.PHONY: service
service: tidy
	@echo "Starting service with $(ENV_FILE)..."
	@if [ -f $(ENV_FILE) ]; then \
		echo "Loading environment from $(ENV_FILE)"; \
		set -a; \
		source $(ENV_FILE); \
		set +a; \
		go run $(MAIN_FILE) $(SERVICE_CMD); \
	else \
		echo "Warning: $(ENV_FILE) not found. Running without environment file."; \
		go run $(MAIN_FILE) $(SERVICE_CMD); \
	fi

# Tidy dependencies
.PHONY: tidy
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

# Run the application
.PHONY: run
run:
	@echo "Running application..."
	@go run $(MAIN_FILE) $(SERVICE_CMD)

# Build the application
.PHONY: build
build:
	@echo "Building application..."
	@go build -o bin/app $(MAIN_FILE)

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@go clean