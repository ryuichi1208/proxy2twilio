# Configuration
HTTP_PORT = 3000
UPSTREAM_URL = http://example.com

# Default target
all: build

# Build the binary
build: $(GO_FILES)
	@echo "Building $(APP_NAME)..."
	@go build -o $(APP_NAME) main.go
	@echo "Build completed."

# Run the application
run: build
	@echo "Running $(APP_NAME) on port $(HTTP_PORT)..."
	@./$(APP_NAME)

# Clean the build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f $(APP_NAME)
	@echo "Cleanup completed."

# Format Go code
fmt:
	@echo "Formatting Go code..."
	@gofmt -w $(GO_FILES)
	@echo "Formatting completed."

# Check for lint issues
lint:
	@echo "Running static analysis..."
	@go vet ./...
	@echo "Static analysis completed."

# Test the application (no tests included here, but you can add them)
test:
	@echo "Running tests..."
	@go test ./...
	@echo "Tests completed."

# Generate a binary release
release: clean build
	@echo "Packaging $(APP_NAME) binary..."
	@tar -czvf $(APP_NAME).tar.gz $(APP_NAME)
	@echo "Release package created: $(APP_NAME).tar.gz"

# Help menu
help:
	@echo "Makefile for $(APP_NAME)"
	@echo
	@echo "Usage:"
	@echo "  make build      - Build the application binary."
	@echo "  make run        - Run the application."
	@echo "  make clean      - Remove build artifacts."
	@echo "  make fmt        - Format Go code."
	@echo "  make lint       - Run static analysis."
	@echo "  make test       - Run tests."
	@echo "  make release    - Build and package the binary for release."
	@echo "  make help       - Show this help message."

.PHONY: all build run clean fmt lint test release help
