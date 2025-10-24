# =============================================================================
# Variables
# =============================================================================

GO_CMD = go

# =============================================================================
# Help
# =============================================================================

.PHONY: help
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# =============================================================================
# Docker Compose
# =============================================================================

.PHONY: infra-up
infra-up: ## Start infrastructure services (MySQL, Kafka, etc.)
	docker compose -f docker-compose.infra.yml up -d

.PHONY: infra-down
infra-down: ## Stop infrastructure services
	docker compose -f docker-compose.infra.yml down

.PHONY: up
up: infra-up ## Build and start application services
	docker compose -f docker-compose.yml up -d --build

.PHONY: down
down: ## Stop application services
	docker compose -f docker-compose.yml down

.PHONY: start
start: up ## Alias for 'up'

.PHONY: stop
stop: down infra-down ## Stop all services (apps and infra)

.PHONY: logs
logs: ## Show logs from all application services
	docker compose -f docker-compose.yml logs -f

.PHONY: logs-starter
logs-starter: ## Show logs from starter-service
	docker compose -f docker-compose.yml logs -f starter-service

.PHONY: logs-notification
logs-notification: ## Show logs from notification-service
	docker compose -f docker-compose.yml logs -f notification-service

.PHONY: rebuild
rebuild: ## Rebuild application services without cache
	docker compose -f docker-compose.yml build --no-cache

# =============================================================================
# Go Development
# =============================================================================

.PHONY: dev-starter
dev-starter: ## Run starter-service locally
	cd services/starter-service && $(GO_CMD) run ./cmd/main.go

.PHONY: dev-notification
dev-notification: ## Run notification-service locally
	cd services/notification-service && $(GO_CMD) run ./cmd/main.go

.PHONY: test
test: ## Run unit tests for all services and packages (excluding integration tests)
	$(GO_CMD) test -short -v ./pkg/...
	cd services/starter-service && $(GO_CMD) test -short -v ./...
	cd services/notification-service && $(GO_CMD) test -short -v ./...

.PHONY: test-integration
test-integration: ## Run integration tests for starter-service (requires Docker)
	@echo "Running integration tests (requires Docker to be running)..."
	cd services/starter-service && $(GO_CMD) test -v ./test/integration/...

.PHONY: test-all
test-all: test test-integration ## Run all tests (unit + integration)

.PHONY: test-coverage
test-coverage: ## Generate test coverage report for starter-service
	cd services/starter-service && $(GO_CMD) test -short -coverprofile=coverage.out -covermode=atomic ./...
	cd services/starter-service && $(GO_CMD) tool cover -html=coverage.out -o coverage.html

.PHONY: test-integration-coverage
test-integration-coverage: ## Generate integration test coverage report
	cd services/starter-service && $(GO_CMD) test -coverprofile=coverage_integration.out -covermode=atomic ./test/integration/...
	cd services/starter-service && $(GO_CMD) tool cover -html=coverage_integration.out -o coverage_integration.html

.PHONY: fmt
fmt: ## Format all Go code
	$(GO_CMD) fmt ./...

.PHONY: lint
lint: ## Lint all Go code
	golangci-lint run ./pkg/...
	golangci-lint run ./services/starter-service/...
	golangci-lint run ./services/notification-service/...

.PHONY: deps
deps: ## Sync and tidy Go modules
	$(GO_CMD) work sync
	$(GO_CMD) mod tidy
	cd pkg && $(GO_CMD) mod tidy
	cd services/starter-service && $(GO_CMD) mod tidy
	cd services/notification-service && $(GO_CMD) mod tidy

# =============================================================================
# Database
# =============================================================================

.PHONY: db-migrate
db-migrate: ## Run database migrations for starter-service
	docker compose -f docker-compose.infra.yml exec mysql mysql -u gouser -pgopassword intern_app < services/starter-service/migrations/001_init.sql
