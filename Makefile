# Generate swagger docs
docs:
	swag init -g cmd/server/main.go -o docs/v1

# Run the application
run:
	go run cmd/server/main.go

# Build and run
start: docs run

# Generate a new module
generate-module:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make generate-module name=<module-name> [description='Module description']"; \
		echo "Example: make generate-module name=garage description='Car garage management'"; \
		exit 1; \
	fi
	@./scripts/generate-module.sh $(name) "$(description)"

# Clean build artifacts
clean:
	rm -rf build/
	rm -rf tmp/
	go clean

# Build the application
build:
	go build -o build/pitstop cmd/server/main.go

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Tidy dependencies
tidy:
	go mod tidy

# Development helpers
dev-setup: tidy fmt
	@echo "Development environment setup complete!"

# Show help
help:
	@echo "Available commands:"
	@echo "  docs             - Generate swagger documentation"
	@echo "  run              - Run the application"
	@echo "  start            - Generate docs and run the application"
	@echo "  generate-module  - Generate a new module (requires name=<module-name>)"
	@echo "  build            - Build the application"
	@echo "  test             - Run tests"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo "  clean            - Clean build artifacts"
	@echo "  fmt              - Format code"
	@echo "  lint             - Run linter"
	@echo "  tidy             - Tidy dependencies"
	@echo "  dev-setup        - Setup development environment"
	@echo "  help             - Show this help message"

.PHONY: docs run start generate-module build test test-coverage clean fmt lint tidy dev-setup help