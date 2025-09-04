.PHONY: build run test clean docker-build docker-run migrate

# Build the application
build:
	go build -o bin/api cmd/api/main.go
	go build -o bin/worker cmd/worker/main.go

# Run the application
run:
	go run cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

# Build Docker image
docker-build:
	docker build -t finfam-backend

# Run with Docker Compose
docker-run:
	docker compose -f docker/docker-compose.yml up -d --build
	docker compose -f docker/docker-compose.yml ps

# Stop Docker Compose
docker-stop:
	docker compose -f docker/docker-compose.yml down

# Database migrations
migrate-up:
	migrate -path migrations -database "postgres://finfam:finfam123@localhost:5434/finfam?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://finfam:finfam123@localhost:5434/finfam?sslmode=disable" down

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Development setup
dev-setup:
	docker-compose up -d postgres redis
	make migrate-up

# Generate API documentation
docs:
	swag init -g cmd/api/main.go