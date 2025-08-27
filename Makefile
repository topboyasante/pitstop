# Generate swagger docs
docs:
	swag init -g cmd/server/main.go -o docs/v1

# Run the application
run:
	go run cmd/server/main.go

# Build and run
start: docs run