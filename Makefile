.PHONY: help build run stop restart logs clean test migrate

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build Docker images
	docker compose build

up: ## Start all services
	docker compose up -d

down: ## Stop all services
	docker compose down

restart: ## Restart all services
	docker compose restart

logs: ## Show logs
	docker compose logs -f

logs-app: ## Show application logs
	docker compose logs -f app

logs-mysql: ## Show MySQL logs
	docker compose logs -f mysql

clean: ## Stop and remove all containers, networks, and volumes
	docker compose down -v

ps: ## Show running containers
	docker compose ps

exec-app: ## Execute shell in app container
	docker compose exec app sh

exec-mysql: ## Execute MySQL shell
	docker compose exec mysql mysql -u gouser -pgopassword intern_app

migrate: ## Run database migrations
	docker compose exec mysql mysql -u gouser -pgopassword intern_app < migrations/001_create_starter_table.sql

rebuild: ## Rebuild and restart
	docker compose down
	docker compose build --no-cache
	docker compose up -d

test: ## Run all tests
	go test ./... -v

test-unit: ## Run only unit tests (fast, no external dependencies)
	@echo "Running unit tests..."
	go test ./internal/.../domain/... -v
	go test ./internal/.../application/... -v -short

test-integration: ## Run integration tests (requires test database)
	@echo "Starting test database..."
	docker-compose -f docker-compose.test.yml up -d
	@echo "Waiting for database to be ready..."
	@sleep 5
	@echo "Running integration tests..."
	go test -tags=integration ./test/integration/... -v || true
	@echo "Stopping test database..."
	docker-compose -f docker-compose.test.yml down

test-e2e: ## Run end-to-end tests (requires full application)
	@echo "Starting test environment..."
	docker-compose -f docker-compose.test.yml up -d
	@sleep 5
	@echo "Running e2e tests..."
	go test -tags=e2e ./test/e2e/... -v || true
	@echo "Stopping test environment..."
	docker-compose -f docker-compose.test.yml down

test-short: ## Run quick tests (skip slow integration/e2e tests)
	go test ./... -short -v

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	go tool cover -func=coverage.out | grep total

test-coverage-unit: ## Run unit tests with coverage
	@echo "Running unit tests with coverage..."
	go test ./internal/.../domain/... ./internal/.../application/... -coverprofile=coverage-unit.out -covermode=atomic
	go tool cover -html=coverage-unit.out -o coverage-unit.html
	@echo "Unit test coverage report generated: coverage-unit.html"
	go tool cover -func=coverage-unit.out | grep total

test-db-up: ## Start test database only
	docker-compose -f docker-compose.test.yml up -d mysql-test
	@echo "Test database is starting on port 3307..."
	@echo "Connection string: mysql://root:test_password@localhost:3307/gorest_test"

test-db-down: ## Stop test database
	docker-compose -f docker-compose.test.yml down

test-db-logs: ## Show test database logs
	docker-compose -f docker-compose.test.yml logs -f mysql-test

clean-test: ## Clean test artifacts and stop test containers
	docker-compose -f docker-compose.test.yml down -v
	rm -f coverage.out coverage.html coverage-unit.out coverage-unit.html

dev: ## Run in development mode (without Docker)
	go run cmd/api/main.go
