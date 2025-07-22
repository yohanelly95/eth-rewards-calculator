# Binary name
BINARY_NAME=eth-rewards
BINARY_DIR=bin

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Build flags
LDFLAGS=-ldflags "-s -w"

# Main package path
MAIN_PATH=./cmd/calculator

all: test build

build:
	@echo "Building..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_DIR)/$(BINARY_NAME)"

build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BINARY_DIR)
	# Linux AMD64
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	# Linux ARM64
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	# Darwin AMD64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	# Darwin ARM64 (M1)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	# Windows AMD64
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "Multi-platform build complete"

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean:
	@echo "Cleaning..."
	@rm -rf $(BINARY_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

vet:
	@echo "Running go vet..."
	$(GOVET) ./...

lint: fmt vet
	@echo "Linting complete"

run: build
	@echo "Running with default parameters..."
	./$(BINARY_DIR)/$(BINARY_NAME) -validators 4096 -detailed

install: build
	@echo "Installing..."
	@cp $(BINARY_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installed to /usr/local/bin/$(BINARY_NAME)"

uninstall:
	@echo "Uninstalling..."
	@rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "Uninstalled"

.PHONY: all build build-all test test-coverage clean deps fmt vet lint run install uninstall