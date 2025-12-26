ENV ?= test
LAMBDA_ENABLED ?= 0

.PHONY: build-app
build-app:
	go build -o bin/api cmd/api/main.go
ifeq ($(LAMBDA_ENABLED), 1)
	zip api.zip bin/api
endif

.PHONY: run
run:
	DB_HOST=localhost SERVER_PORT=8080  air -c .air.toml

.PHONY: into_db
into_db:
	docker compose -f docker-compose.yml exec -it kinance_db psql -U finfam -d finfam

# Run tests
test:
	go test -v ./...

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
	docker-compose down
	docker-compose up -d

.PHONY: docs
docs:
	swag init -g cmd/api/main.go

.PHONY: redis-cli
redis-cli:
	docker-compose exec -it kinance_redis redis-cli -h localhost -p 6379

.PHONY: tf-init
ifeq ($(ENV), test)
tf-init:
	tflocal -chdir=infra/env/test/ init
else
tf-init:
	@echo "Not ready ..."
endif

.PHONY: tf-plan
ifeq ($(ENV), test) 
tf-plan:
	tflocal -chdir=infra/env/test/ plan
else
tf-plan:
	@echo "Not ready ..."
endif

.PHONY: tf-apply
ifeq ($(ENV), test) 
tf-apply:
	tflocal -chdir=infra/env/test/ apply -auto-approve
else
tf-apply:
	@echo "Not ready ..."
endif

.PHONY: lstack-run
lstack-run:
	localstack 