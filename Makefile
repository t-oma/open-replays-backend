server:
	@echo "Starting server..."
	go run ./cmd/api/main.go

test:
	@echo "Running tests..."
	go test -v ./...
