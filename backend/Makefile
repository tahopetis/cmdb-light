.PHONY: test test-coverage test-unit test-integration clean build run deps fmt lint vet

# Default target
all: deps test build

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Run vet
vet:
	go vet ./...

# Run unit tests
test-unit:
	go test -v -short ./...

# Run integration tests
test-integration:
	go test -v -run=Integration ./...

# Run all tests
test: test-unit

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -f coverage.out coverage.html
	rm -f cmdb-lite-backend

# Build the application
build:
	go build -o cmdb-lite-backend ./cmd/main.go

# Run the application
run:
	go run ./cmd/main.go

# Run tests and check coverage
check: test-coverage

# Format, lint, and test
ci: fmt lint vet test-coverage

# Run the application in development mode
dev:
	go run ./cmd/main.go --config=./configs/config.dev.yml