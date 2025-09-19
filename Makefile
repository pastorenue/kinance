.PHONY: build-app run-app test clean docker-build docker-run migrate

.PHONY: build-app
build-app:
	go build -o bin/api cmd/api/main.go
	go build -o bin/worker cmd/worker/main.go

.PHONY: run-app
run-app:
	go run cmd/api/main.go

.PHONY: into_db
into_db:
	docker compose -f docker-compose.yml exec -it kinance_db psql -U finfam -d finfam

# Run tests
test:
	go test -v -cover ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: clean
clean:
	rm -rf bin/
	rm -f coverage.out

.PHONY: docker-build
docker-build:
	docker build -t finfam-backend .

.PHONY: docker-run
docker-run:
	docker compose -f docker-compose.yml up -d --build
	docker compose -f docker-compose.yml ps

.PHONY: docker-stop
docker-stop:
	docker compose -f docker-compose.yml down

.PHONY: docker-restart
docker-restart:
	make docker-stop
	make docker-run

.PHONY: migrate
migrate-up:
	migrate -path migrations -database "postgres://finfam:finfam123@localhost:5434/finfam?sslmode=disable" up

.PHONY: migrate-down
migrate-down:
	migrate -path migrations -database "postgres://finfam:finfam123@localhost:5434/finfam?sslmode=disable" down

.PHONY: deps
deps:
	go mod download
	go mod tidy

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: dev-setup
dev-setup:
	docker-compose up -d postgres redis
	make migrate-up

.PHONY: docs
docs:
	swag init -g cmd/api/main.go
