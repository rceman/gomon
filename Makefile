# Makefile for gomon project

# Define variables
BINARY_NAME=gomon
BUILD_DIR=bin

# Default target
.PHONY: all
all: build

# Build target
.PHONY: build
build:
	@echo "Building the project..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/$(BINARY_NAME)

# Clean target
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -f $(BUILD_DIR)/$(BINARY_NAME)

# Run target
.PHONY: run
run: build
	@echo "Running the project..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Test target
.PHONY: test
test:
	@echo "Running tests..."
	@go test ./...

# Help target
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make                   Build the project (default)"
	@echo "  make build             Build the project"
	@echo "  make clean             Clean the build artifacts"
	@echo "  make run               Build and run the project"
	@echo "  make test              Run the tests"
	@echo "  make help              Show this help message"
